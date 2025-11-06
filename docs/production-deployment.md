# Production Database Deployment Guide

This guide covers how to set up your database on a production server that doesn't have your project code or Docker.

## Overview

There are **3 main approaches** to deploy your database to production:

1. **Migrations Only** (Recommended) - Run migrations on empty production DB
2. **Data Migration** - Copy dev data to production (for initial seeding)
3. **Hybrid** - Run migrations + seed production-appropriate data

---

## Option 1: Migrations Only (Recommended for Production)

This is the **cleanest approach** for production. You only need the migration files and the `migrate` tool.

### On Production Server:

```bash
# 1. Install golang-migrate
# For Linux:
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/migrate

# For macOS:
brew install golang-migrate

# 2. Copy only migration files to server
# From your local machine:
scp -r db/schema user@production-server:/path/to/migrations/

# 3. Run migrations on production server
migrate -path /path/to/migrations \
        -database "postgres://user:password@localhost:5432/production_db?sslmode=disable" \
        up

# 4. (Optional) Seed with production data
# Either use the Go seeder or manually insert data
```

**Pros:**
- ✅ Clean, reproducible
- ✅ No dev data in production
- ✅ Standard practice
- ✅ Only need migration files

**Cons:**
- ❌ Database starts empty
- ❌ Need to seed production data separately

---

## Option 2: Data Migration (Copy Dev Data)

Use this if you need to copy your dev database (with data) to production for initial setup.

### Step 1: Backup Dev Database

```bash
# From your local machine
make db-backup

# Or manually:
./scripts/backup-db.sh backups/production-init.sql
```

### Step 2: Transfer to Production

```bash
# Copy backup file to production server
scp backups/production-init.sql user@production-server:/tmp/

# Or use a secure file transfer method
```

### Step 3: Restore on Production

On the production server:

```bash
# Method A: Using psql directly
psql -U postgres -d production_db < /tmp/production-init.sql

# Method B: If you copied the restore script
./scripts/restore-db.sh /tmp/production-init.sql \
  "postgres://user:password@localhost:5432/production_db"
```

**Pros:**
- ✅ Quick initial setup
- ✅ Includes all dev data
- ✅ Good for testing/staging

**Cons:**
- ❌ Dev data in production (not ideal)
- ❌ Larger file transfer
- ❌ May contain sensitive dev data

---

## Option 3: Hybrid Approach (Migrations + Production Seeding)

Best of both worlds: clean migrations + production-appropriate data.

### Step 1: Setup Database Structure

```bash
# On production server
migrate -path /path/to/migrations \
        -database "postgres://user:password@localhost:5432/production_db" \
        up
```

### Step 2: Seed Production Data

**Option A: Upload and run seeder binary**

```bash
# On your local machine - build the seeder
GOOS=linux GOARCH=amd64 go build -o seeder-linux cmd/seeder/main.go

# Transfer to production
scp seeder-linux user@production-server:/tmp/

# On production server
DATABASE_URL="postgres://user:password@localhost/production_db" \
  /tmp/seeder-linux -users=1000 -orders=5000
```

**Option B: SQL seed file**

```bash
# Create production seed file locally
cat > production-seed.sql << 'EOF'
-- Insert admin user
INSERT INTO users (email, username, password_hash, first_name, last_name)
VALUES ('admin@yourcompany.com', 'admin', '$2a$...', 'Admin', 'User');

-- Add other production-specific data
EOF

# Transfer and run
scp production-seed.sql user@production-server:/tmp/
psql -U user -d production_db < /tmp/production-seed.sql
```

**Pros:**
- ✅ Clean migrations
- ✅ Production-appropriate data
- ✅ Reproducible
- ✅ No dev artifacts

**Cons:**
- ❌ More steps
- ❌ Need to maintain seed data

---

## Quick Reference: What to Transfer

### Minimal (Migrations Only)
```
db/schema/
  ├── 001_create_users_table.up.sql
  ├── 001_create_users_table.down.sql
  ├── 002_create_orders_table.up.sql
  └── 002_create_orders_table.down.sql
```

### With Seeder
```
db/schema/          # Migration files
seeder-linux        # Compiled seeder binary
```

### Full Backup
```
backups/production-init.sql   # Complete database dump
```

---

## Production Checklist

- [ ] Production PostgreSQL server is running
- [ ] Database user and database created
- [ ] Firewall allows connections (if remote)
- [ ] SSL/TLS configured (recommended)
- [ ] Backup strategy in place
- [ ] Environment variables set
- [ ] Test migration rollback works
- [ ] Document database credentials securely

---

## Environment Variables for Production

Create a `.env` file or set system variables:

```bash
DATABASE_URL="postgres://username:password@host:5432/dbname?sslmode=require"
PORT="8080"
```

**Security Notes:**
- ✅ Use strong passwords
- ✅ Enable SSL (`sslmode=require`)
- ✅ Restrict database user permissions
- ✅ Use connection pooling
- ✅ Never commit `.env` to git

---

## Common Production Setup Commands

### Create Production Database
```bash
# On production PostgreSQL server
sudo -u postgres psql

CREATE DATABASE production_db;
CREATE USER app_user WITH PASSWORD 'strong_password';
GRANT ALL PRIVILEGES ON DATABASE production_db TO app_user;

# Enable UUID extension
\c production_db
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
```

### Run Migrations
```bash
migrate -path db/schema \
        -database "postgres://app_user:strong_password@localhost:5432/production_db?sslmode=require" \
        up
```

### Verify Migration Status
```bash
migrate -path db/schema \
        -database "postgres://app_user:strong_password@localhost:5432/production_db?sslmode=require" \
        version
```

---

## Troubleshooting

### "migrate: command not found"
```bash
# Install golang-migrate
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/
```

### "uuid_generate_v4() does not exist"
```bash
# Enable UUID extension
psql -U postgres -d production_db -c 'CREATE EXTENSION IF NOT EXISTS "uuid-ossp";'
```

### "permission denied"
```bash
# Grant proper permissions
psql -U postgres -c "GRANT ALL PRIVILEGES ON DATABASE production_db TO app_user;"
```

---

## Recommended Production Workflow

1. **First Time Setup:**
   - Transfer migration files
   - Run migrations on empty production DB
   - Create initial admin user manually or via seed script

2. **Future Updates:**
   - Create new migration files in dev
   - Test migrations locally
   - Transfer new migration files to production
   - Run `migrate up` on production

3. **Data Management:**
   - Use application to create production data
   - Or use seeder with production-appropriate values
   - Keep backups before major changes

---

## Additional Resources

- [golang-migrate documentation](https://github.com/golang-migrate/migrate)
- [PostgreSQL backup/restore](https://www.postgresql.org/docs/current/backup.html)
- See `docs/migrate-behavior-and-safety.md` for migration details
