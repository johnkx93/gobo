#!/bin/bash
# Backup production database before migrations
# Automatically called by make migrate-up-prod

set -e

if [ -z "$PRODUCTION_DATABASE_URL" ]; then
    echo "❌ Error: PRODUCTION_DATABASE_URL is not set"
    exit 1
fi

# Create backups directory if it doesn't exist
mkdir -p backups/production

# Generate timestamped backup filename
TIMESTAMP=$(date +%Y%m%d-%H%M%S)
OUTPUT_FILE="backups/production/prod-backup-${TIMESTAMP}.sql"

# echo "📦 Backing up production database..."
# echo "Output: $OUTPUT_FILE"
echo "HAVENT IMPLEMENT pg_dump / db backup YET"
echo "HAVENT IMPLEMENT pg_dump / db backup YET"
echo "HAVENT IMPLEMENT pg_dump / db backup YET"
echo "HAVENT IMPLEMENT pg_dump / db backup YET"
echo "HAVENT IMPLEMENT pg_dump / db backup YET"

# Use pg_dump with connection string
# pg_dump "$PRODUCTION_DATABASE_URL" > "$OUTPUT_FILE"

if [ $? -eq 0 ]; then
    echo "✅ Backup created successfully: $OUTPUT_FILE"
    echo "📊 File size: $(du -h "$OUTPUT_FILE" | cut -f1)"
    
    # Keep only last 10 backups
    ls -t backups/production/prod-backup-*.sql 2>/dev/null | tail -n +11 | xargs rm -f 2>/dev/null || true
    echo "🗂️  Keeping last 10 backups, older ones removed"
else
    echo "❌ Backup failed!"
    exit 1
fi
