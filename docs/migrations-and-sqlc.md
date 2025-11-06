# Migrations + sqlc Guide

This guide explains the recommended workflow for making schema changes, updating SQL query files, and regenerating the typed sqlc client in this project.

Prerequisites
- Docker & docker-compose (for local Postgres used by the project)
- `migrate` CLI (golang-migrate)
- `sqlc` CLI
- Go toolchain (>= 1.21)

Key paths
- Migrations: `db/schema/` (files created by `make migrate-create`)
- SQL queries used by sqlc: `db/queries/` (e.g., `db/queries/user.sql`)
- Generated sqlc client: `internal/db/` (DO NOT edit these files manually)
- Docker DB init scripts: `docker/postgres/` (run only on first container init)

High-level workflow
1. Create migration(s) to change the schema (if needed).
2. Edit `db/queries/*.sql` to add/modify queries that depend on the new schema.
3. Apply migrations to the local DB.
4. Run `sqlc generate` to update the typed client.
5. Update Go DTOs / handlers / services to use new/generated types.
6. Run tests and verify the app.

Step-by-step

1) Start local Postgres (if using Dockerized DB supplied with this repo)

```bash
make docker-up
```

This runs the `postgres:16-alpine` container configured in `docker-compose.yml`. The DB is exposed on `localhost:5432`.

2) Create a new migration skeleton

Use the Makefile helper to create a pair of migration files:

```bash
make migrate-create name=add_phone_to_users
```

This creates two files under `db/schema/` like:
- `db/schema/00X_add_phone_to_users.up.sql`
- `db/schema/00X_add_phone_to_users.down.sql`

Edit the `.up.sql` to apply your change (example: add a `phone` column):

```sql
ALTER TABLE users
ADD COLUMN phone TEXT;
```

Edit the `.down.sql` to roll it back:

```sql
ALTER TABLE users
DROP COLUMN phone;
```

Notes:
- Use `IF NOT EXISTS` or similar defensive SQL if you expect the migration might run against variable states during development.
- Migrations are the canonical way to change schema; do not modify `docker/postgres/init.sql` for repeatable schema changes (that file runs only on first DB init).

3) Modify `db/queries/user.sql`

Update the queries to reflect the new schema. For example, if you add a `phone` column, include it in INSERT/SELECT/UPDATE queries.

Example additions/updates for sqlc:

```sql
-- name: CreateUser :one
INSERT INTO users (id, email, password_hash, phone)
VALUES (:id, :email, :password_hash, :phone)
RETURNING id, email, phone, created_at;

-- name: GetUserByID :one
SELECT id, email, phone, created_at
FROM users
WHERE id = $1;
```

Guidelines for sqlc queries:
- Use explicit column lists in SELECT/INSERT/UPDATE.
- Add a `-- name: <FuncName> :one` or `:many` comment above statements to control generation.
- Keep placeholders consistent with sqlc expectations (use numbered placeholders like $1 for Postgres or named parameters depending on your sqlc config).

4) Apply migrations locally

Run:

```bash
make migrate-up
```

This uses the `migrate` CLI and the connection string configured in the Makefile (defaults to `postgres://postgres:postgres@localhost:5432/appdb?sslmode=disable`). Ensure the DB is up (`make docker-up`) before running migrations.

If you need to rollback a migration:

```bash
make migrate-down
```

If you want to force a specific migration version (use carefully):

```bash
make migrate-force version=5
```

5) Regenerate sqlc typed client

Run:

```bash
make sqlc-generate
```

This executes `sqlc generate` using the repository's `sqlc.yaml` and updates the generated code under `internal/db/` (or the configured output directory).

Important: do not edit generated files under `internal/db/` directly. Make changes in `db/queries/*.sql` and re-run `sqlc generate`.

6) Update Go code and tests

- Update DTOs in `internal/app/user/dto.go` to include new fields and validation tags.
- Update the service layer `internal/app/user/service.go` to map from sqlc-generated types to your domain types (and vice versa).
- Update handlers `internal/app/user/handler.go` to decode the new fields from requests and encode them in responses.
- Update unit/integration tests to reflect the new schema and query outputs.

Run tests and fix compile/runtime errors:

```bash
make test
# or
go test ./...
```

7) Run and verify

Start the application:

```bash
make run
```

Exercise endpoints (curl, Postman) and verify DB rows:

```bash
# example: create user (adjust path/payload to your API)
curl -X POST http://localhost:8080/api/v1/users \
  -H 'Content-Type: application/json' \
  -d '{"email":"alice@example.com","password":"secret","phone":"+1-555-0123"}'
```

Debugging and troubleshooting

- If `sqlc generate` errors:
  - Check SQL syntax in `db/queries/*.sql`.
  - Inspect `sqlc.yaml` for correct input/output paths and Postgres dialect.
- If migrations fail:
  - Ensure DB is up and listening on the expected port.
  - If migration says column already exists, either update the migration to use `IF NOT EXISTS` or reset the dev DB.
- If queries compile but your Go code fails to build:
  - Update imports and struct mappings to use the new generated types.

Resetting DB (to re-run `docker/postgres` init scripts)

The SQL files in `docker/postgres/` are executed only when the Postgres container initializes a new data directory. To force full reinitialization (DESTROYS DATA):

```bash
docker-compose down -v
make docker-up
make migrate-up
```

Be careful: `docker-compose down -v` deletes all named volumes, including `postgres_data`.

Quick checklist

- [ ] Create migration via `make migrate-create name=...`
- [ ] Edit `db/schema/*.up.sql` and `*.down.sql`
- [ ] Update `db/queries/*.sql` for new queries
- [ ] Run `make docker-up` (if needed)
- [ ] Run `make migrate-up`
- [ ] Run `make sqlc-generate`
- [ ] Update Go layers (DTOs, services, handlers)
- [ ] Run `make test` and `make run`

Where to look for help
- `db/schema/` — migrations
- `db/queries/` — sqlc source SQL
- `internal/db/` — generated code (read-only)
- `docker/postgres/init.sql` — container init SQL (runs once on new volume)
- Makefile targets: `docker-up`, `migrate-up`, `sqlc-generate`, `run`, `setup`, `dev`