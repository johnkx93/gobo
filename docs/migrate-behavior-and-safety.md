# Migration behavior and safety notes

This document expands on migration behavior (golang-migrate) and practical safety guidance for this project.

Contents
- How versions are determined
- What `migrate up` does
- What `migrate down` does
- Examples and common commands
- What happens on failure
- Practical recommendations and safety checklist
- Resetting DB and init scripts

---

## How versions are determined

Migrations are ordered by the version prefix in the migration filenames (e.g. `003_add_thing.up.sql`, `004_do_more.up.sql`, or timestamp-based prefixes). The migration tool applies `up` files in ascending order and `down` files in descending order.

In this repo the migration files live in `db/schema/`.

## What `migrate up` does

- `migrate up` applies all unapplied `*.up.sql` migration files in ascending version order.
- If the database's current recorded version is `002` and you have `003` through `010` present, `migrate up` will attempt to apply `003`, then `004`, ... up to `010` sequentially (it does not stop after applying only `003`).

## What `migrate down` does

- `migrate down` with no numeric argument reverts all applied migrations (rolls back to version 0).
- `migrate down <N>` reverts the last N applied migrations (some `migrate` CLI variants support `down N`).
- In this repository `make migrate-down` runs `migrate ... down 1`, so the Makefile target will revert only the most recently applied migration (one step).

You can also use `migrate goto <version>` to move the database to a specific version (this runs up or down as necessary).

## Examples and common commands

### Development (Local Docker Database)

Using the connection string from the Makefile:

- Show current database migration version:

```bash
migrate -path db/schema -database "postgres://postgres:postgres@localhost:5432/appdb?sslmode=disable" version
```

- Apply all pending migrations:

```bash
migrate -path db/schema -database "postgres://postgres:postgres@localhost:5432/appdb?sslmode=disable" up
# or with Makefile
make migrate-up
```

- Revert the last applied migration (one step):

```bash
migrate -path db/schema -database "postgres://postgres:postgres@localhost:5432/appdb?sslmode=disable" down 1
# or with Makefile
make migrate-down
```

- Revert 3 migrations:

```bash
migrate -path db/schema -database "postgres://postgres:postgres@localhost:5432/appdb?sslmode=disable" down 3
```

- Revert everything (to version 0):

```bash
migrate -path db/schema -database "postgres://postgres:postgres@localhost:5432/appdb?sslmode=disable" down
```

- Force-set recorded version (dangerous; use only for recovery):

```bash
migrate -path db/schema -database "..." force <version>
# or with Makefile
make migrate-force version=3
```

- Move to a specific version:

```bash
migrate -path db/schema -database "..." goto <version>
```

### Production Database

**Set up environment variable first:**
```bash
export PRODUCTION_DATABASE_URL='postgres://user:pass@host:5432/dbname?sslmode=require'
```

**Available production commands:**

- Check migration status:
```bash
make migrate-status-prod
```

- Apply all pending migrations (with automatic backup):
```bash
make migrate-up-prod
# This automatically:
#   1. Creates backup in backups/production/
#   2. Shows target database (password hidden)
#   3. Asks for confirmation
#   4. Runs migrations
```

- Rollback last migration (use with extreme caution):
```bash
make migrate-down-prod
# Includes confirmation prompt
```

**Manual production commands (if needed):**

```bash
# Check version
migrate -path db/schema -database "$PRODUCTION_DATABASE_URL" version

# Apply migrations
migrate -path db/schema -database "$PRODUCTION_DATABASE_URL" up

# Rollback one migration
migrate -path db/schema -database "$PRODUCTION_DATABASE_URL" down 1
```

> Note: exact CLI flags and supported subcommands depend on your installed `migrate` binary version; check `migrate --help` if anything differs.

## What happens on failure

- Each migration file (an `up` or `down` script) is executed as a single unit. For Postgres, the tool runs each migration inside a transaction (where supported), so if a migration fails its changes in that migration are rolled back.
- If a migration fails partway (for example migration `005` fails while `003` and `004` have already been applied), `003` and `004` remain applied and the CLI stops on the failing migration. You'll need to fix `005` and run `migrate up` again.
- If you run `down N` and it fails on the second step, the first step already reverted remains reverted; recovery requires fixing the failing down migration or using `force`/`goto` to correct recorded version state.

## Practical recommendations and safety checklist

### Development safety
- Prefer small, incremental migrations to reduce blast radius.
- Apply migrations locally first (dockerized DB or a disposable test DB) and verify before applying to shared environments.
- Use `make migrate-down` (the repo's `down 1`) to roll back the last migration during development.

### Production safety

**CRITICAL: Always backup before production migrations!**

This project includes automatic backup protection:

- `make migrate-up-prod` automatically creates a timestamped backup before running migrations
- Backups are stored in `backups/production/` with timestamps
- Last 10 backups are kept automatically; older ones are cleaned up
- Manual backups can be created with `make db-backup`

**Production Migration Workflow:**

```bash
# 1. Set production database URL
export PRODUCTION_DATABASE_URL='postgres://user:pass@host:5432/db?sslmode=require'

# 2. Check current migration status
make migrate-status-prod

# 3. Run migrations (automatic backup + confirmation prompt)
make migrate-up-prod
# This will:
#   - Create automatic backup to backups/production/
#   - Show which database you're targeting
#   - Ask for confirmation
#   - Apply migrations
```

**If something goes wrong:**

```bash
# Restore from automatic backup
make db-restore file=backups/production/prod-backup-YYYYMMDD-HHMMSS.sql

# Or manually with psql
psql "$PRODUCTION_DATABASE_URL" < backups/production/prod-backup-YYYYMMDD-HHMMSS.sql
```

**Additional production safety rules:**

- Never run `down` (full rollback) on production without a tested recovery plan.
- Test rollback migrations (`down` scripts) in staging before production use.
- Consider writing reversible migrations (pairs of up/down) and avoid data-loss operations without migration-specific backup/restore steps.
- Monitor migration execution time on large tables; some migrations may require maintenance windows.
- Use transactions where possible (Postgres default) to ensure atomic migrations.

### Backup strategy

**Development backups:**
```bash
make db-backup                    # Manual backup
make db-backup file=backups/before-feature.sql
```

**Production backups:**
```bash
# Automatic: happens before every migrate-up-prod
make migrate-up-prod

# Manual production backup (if needed)
pg_dump "$PRODUCTION_DATABASE_URL" > backups/manual-prod-backup.sql
```

**Restore:**
```bash
# Development (uses Docker)
make db-restore file=backups/backup.sql

# Production (direct connection)
psql "$PRODUCTION_DATABASE_URL" < backups/production/prod-backup-*.sql
```

### Inspect before running
  - Compare the files in `db/schema/` with the output of `migrate version` to see what is pending.
  - If you only want to apply a subset of migrations, either use `up <steps>` (if supported) or create/organize migrations accordingly.

- Recovery and tricky states
  - If the recorded migration version gets out of sync with actual schema (due to manual changes), `migrate force <version>` can set the recorded version but does not alter schema — use it only with full knowledge of schema state.
  - `migrate goto <version>` can be used to move to a specific version but behaves like running multiple up/down steps, so use with caution.

## Resetting DB and init scripts

- Docker init scripts in `docker/postgres/` (for example `docker/postgres/init.sql`) run only when the Postgres container is first created and the data directory is empty.
- To force those init scripts to run again during development (this deletes the persisted DB data):

```bash
# WARNING: deletes DB data stored in named volumes
docker-compose down -v
make docker-up
make migrate-up
```

## Quick safety checklist

### Development
- [ ] Confirm DB is up (`make docker-up`)
- [ ] Inspect `db/schema/` files and current `migrate version`
- [ ] Run `make migrate-up` to apply pending migrations (applies all pending in order)
- [ ] If rollback is needed, prefer `make migrate-down` (reverts last migration) or `down N` for controlled steps
- [ ] Create manual backup before risky operations: `make db-backup`

### Production
- [ ] Set `PRODUCTION_DATABASE_URL` environment variable
- [ ] Check migration status: `make migrate-status-prod`
- [ ] Verify you're targeting the correct database
- [ ] Run `make migrate-up-prod` (automatic backup included)
- [ ] Confirm when prompted
- [ ] Monitor application after migration
- [ ] Keep backup location handy in case of rollback: `backups/production/`
- [ ] Test in staging environment first if possible

### Emergency Rollback
If production migration fails or causes issues:

1. **Restore from automatic backup:**
   ```bash
   # Find the latest backup
   ls -lt backups/production/
   
   # Restore it
   psql "$PRODUCTION_DATABASE_URL" < backups/production/prod-backup-YYYYMMDD-HHMMSS.sql
   ```

2. **Or use migration rollback (if down migration exists and is tested):**
   ```bash
   make migrate-down-prod
   ```

3. **Fix the issue locally, then redeploy:**
   ```bash
   # Fix migration locally
   make migrate-down
   # Edit migration file
   make migrate-up
   # Test thoroughly
   make migrate-up-prod  # Deploy fixed version
   ```

---