#!/bin/bash
# Backup database from Docker container
# Usage: ./scripts/backup-db.sh [output-file]

set -e

# Default output file
OUTPUT_FILE="${1:-backups/db-backup-$(date +%Y%m%d-%H%M%S).sql}"

# Create backups directory if it doesn't exist
mkdir -p "$(dirname "$OUTPUT_FILE")"

echo "📦 Creating database backup..."
echo "Output: $OUTPUT_FILE"

# Backup using docker exec
docker exec -t $(docker-compose ps -q postgres) pg_dump -U postgres appdb > "$OUTPUT_FILE"

echo "✅ Backup created successfully: $OUTPUT_FILE"
echo "📊 File size: $(du -h "$OUTPUT_FILE" | cut -f1)"
