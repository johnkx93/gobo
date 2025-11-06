#!/bin/bash
# Restore database from backup file
# Usage: ./scripts/restore-db.sh <backup-file> <database-url>

set -e

if [ $# -lt 1 ]; then
    echo "Usage: $0 <backup-file> [database-url]"
    echo "Example: $0 backups/db-backup.sql"
    echo "Example: $0 backups/db-backup.sql postgres://user:pass@prod-server/dbname"
    exit 1
fi

BACKUP_FILE="$1"
DB_URL="${2:-postgres://postgres:postgres@localhost:5432/appdb?sslmode=disable}"

if [ ! -f "$BACKUP_FILE" ]; then
    echo "❌ Error: Backup file not found: $BACKUP_FILE"
    exit 1
fi

echo "🔄 Restoring database from: $BACKUP_FILE"
echo "Target: $DB_URL"
echo ""
read -p "⚠️  This will REPLACE all data in the target database. Continue? (yes/no): " confirm

if [ "$confirm" != "yes" ]; then
    echo "❌ Restore cancelled"
    exit 0
fi

echo "📥 Restoring database..."
psql "$DB_URL" < "$BACKUP_FILE"

echo "✅ Database restored successfully!"
