### 1. Check Existing Partitions
-- List all audit_logs partitions (MONTHLY)
SELECT 
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
FROM pg_tables
WHERE tablename LIKE 'audit_logs_%'
ORDER BY tablename;

-- List all error_logs partitions (MONTHLY)
SELECT 
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
FROM pg_tables
WHERE tablename LIKE 'error_logs_%'
ORDER BY tablename;

-- List all orders partitions (YEARLY)
SELECT 
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
FROM pg_tables
WHERE tablename LIKE 'orders_%'
ORDER BY tablename;

-- Check partition details and constraints
SELECT 
    parent.relname AS parent_table,
    child.relname AS partition_name,
    pg_get_expr(child.relpartbound, child.oid) AS partition_range
FROM pg_inherits
JOIN pg_class parent ON pg_inherits.inhparent = parent.oid
JOIN pg_class child ON pg_inherits.inhrelid = child.oid
WHERE parent.relname IN ('audit_logs', 'error_logs', 'orders')
ORDER BY parent.relname, child.relname;

### 2. Create New Partitions

#### For MONTHLY partitions (audit_logs, error_logs):
# For the next month (e.g., January 2026):
-- Create audit_logs partition for January 2026
CREATE TABLE audit_logs_2026_01 PARTITION OF audit_logs
    FOR VALUES FROM ('2026-01-01 00:00:00+00') TO ('2026-02-01 00:00:00+00');

-- Create error_logs partition for January 2026
CREATE TABLE error_logs_2026_01 PARTITION OF error_logs
    FOR VALUES FROM ('2026-01-01 00:00:00+00') TO ('2026-02-01 00:00:00+00');

# For multiple months at once:
-- February 2026
CREATE TABLE audit_logs_2026_02 PARTITION OF audit_logs
    FOR VALUES FROM ('2026-02-01 00:00:00+00') TO ('2026-03-01 00:00:00+00');

CREATE TABLE error_logs_2026_02 PARTITION OF error_logs
    FOR VALUES FROM ('2026-02-01 00:00:00+00') TO ('2026-03-01 00:00:00+00');

-- March 2026
CREATE TABLE audit_logs_2026_03 PARTITION OF audit_logs
    FOR VALUES FROM ('2026-03-01 00:00:00+00') TO ('2026-04-01 00:00:00+00');

CREATE TABLE error_logs_2026_03 PARTITION OF error_logs
    FOR VALUES FROM ('2026-03-01 00:00:00+00') TO ('2026-04-01 00:00:00+00');

#### For YEARLY partitions (orders):
# For the next year (e.g., 2027):
-- Create orders partition for 2027
CREATE TABLE orders_2027 PARTITION OF orders
    FOR VALUES FROM ('2027-01-01 00:00:00+00') TO ('2028-01-01 00:00:00+00');

# For multiple years at once:
-- Create orders partitions for 2028 and 2029
CREATE TABLE orders_2028 PARTITION OF orders
    FOR VALUES FROM ('2028-01-01 00:00:00+00') TO ('2029-01-01 00:00:00+00');

CREATE TABLE orders_2029 PARTITION OF orders
    FOR VALUES FROM ('2029-01-01 00:00:00+00') TO ('2030-01-01 00:00:00+00');

### 3. Check Data Distribution
-- Count records in each audit_logs partition
SELECT 
    tableoid::regclass AS partition_name,
    COUNT(*) AS row_count,
    MIN(created_at) AS oldest_record,
    MAX(created_at) AS newest_record
FROM audit_logs
GROUP BY tableoid
ORDER BY partition_name;

-- Count records in each error_logs partition
SELECT 
    tableoid::regclass AS partition_name,
    COUNT(*) AS row_count,
    MIN(created_at) AS oldest_record,
    MAX(created_at) AS newest_record
FROM error_logs
GROUP BY tableoid
ORDER BY partition_name;

-- Count records in each orders partition
SELECT 
    tableoid::regclass AS partition_name,
    COUNT(*) AS row_count,
    MIN(created_at) AS oldest_record,
    MAX(created_at) AS newest_record
FROM orders
GROUP BY tableoid
ORDER BY partition_name;


### 4. Archive Old Partitions

#### For MONTHLY partitions (audit_logs, error_logs):
# Step 1: Check what you want to archive
-- See partitions older than 1 year
SELECT 
    tablename,
    pg_size_pretty(pg_total_relation_size('public.'||tablename)) AS size
FROM pg_tables
WHERE tablename LIKE 'audit_logs_2024%'
   OR tablename LIKE 'error_logs_2024%'
ORDER BY tablename;

# Step 2: Detach the partition (makes it a standalone table)
-- Detach partition (PostgreSQL 14+)
ALTER TABLE audit_logs DETACH PARTITION audit_logs_2024_01 CONCURRENTLY;

-- For older PostgreSQL versions (without CONCURRENTLY)
ALTER TABLE audit_logs DETACH PARTITION audit_logs_2024_01;

# Step 3: Optional - Backup before dropping
# From terminal
docker exec -t $(docker-compose ps -q postgres) \
    pg_dump -U postgres -t audit_logs_2024_01 appdb > backups/audit_logs_2024_01.sql

# Complete example for archiving:
-- Archive audit_logs from January 2024
ALTER TABLE audit_logs DETACH PARTITION audit_logs_2024_01 CONCURRENTLY;
DROP TABLE IF EXISTS audit_logs_2024_01;

-- Archive error_logs from January 2024
ALTER TABLE error_logs DETACH PARTITION error_logs_2024_01 CONCURRENTLY;
DROP TABLE IF EXISTS error_logs_2024_01;

#### For YEARLY partitions (orders):
# Step 1: Check what you want to archive
-- See orders partitions older than 3 years
SELECT 
    tablename,
    pg_size_pretty(pg_total_relation_size('public.'||tablename)) AS size
FROM pg_tables
WHERE tablename LIKE 'orders_2022'
   OR tablename LIKE 'orders_2021'
ORDER BY tablename;

# Step 2: Detach the partition (makes it a standalone table)
-- Detach partition (PostgreSQL 14+)
ALTER TABLE orders DETACH PARTITION orders_2022 CONCURRENTLY;

# Step 3: Optional - Backup before dropping
# From terminal
docker exec -t $(docker-compose ps -q postgres) \
    pg_dump -U postgres -t orders_2022 appdb > backups/orders_2022.sql

# Step 4: Drop the old partition
DROP TABLE IF EXISTS orders_2022;

### 5. Quick Commands via psql
# Connect to database:
# Via Docker
docker exec -it $(docker-compose ps -q postgres) psql -U postgres -d appdb

# Direct connection
psql -U postgres -d appdb

# Create MONTHLY partitions for next 6 months (audit_logs, error_logs):
-- Copy and paste this in psql
DO $$
DECLARE
    start_date DATE;
    end_date DATE;
    partition_name TEXT;
    year_month TEXT;
BEGIN
    FOR i IN 1..6 LOOP
        start_date := DATE_TRUNC('month', CURRENT_DATE + (i || ' months')::INTERVAL);
        end_date := start_date + INTERVAL '1 month';
        year_month := TO_CHAR(start_date, 'YYYY_MM');
        
        -- Create audit_logs partition
        partition_name := 'audit_logs_' || year_month;
        EXECUTE format(
            'CREATE TABLE IF NOT EXISTS %I PARTITION OF audit_logs FOR VALUES FROM (%L) TO (%L)',
            partition_name, start_date, end_date
        );
        RAISE NOTICE 'Created partition: %', partition_name;
        
        -- Create error_logs partition
        partition_name := 'error_logs_' || year_month;
        EXECUTE format(
            'CREATE TABLE IF NOT EXISTS %I PARTITION OF error_logs FOR VALUES FROM (%L) TO (%L)',
            partition_name, start_date, end_date
        );
        RAISE NOTICE 'Created partition: %', partition_name;
    END LOOP;
END $$;

# Create YEARLY partitions for next 3 years (orders):
-- Copy and paste this in psql
DO $$
DECLARE
    start_date DATE;
    end_date DATE;
    partition_name TEXT;
    year_val TEXT;
BEGIN
    FOR i IN 1..3 LOOP
        start_date := DATE_TRUNC('year', CURRENT_DATE + (i || ' years')::INTERVAL);
        end_date := start_date + INTERVAL '1 year';
        year_val := TO_CHAR(start_date, 'YYYY');
        
        -- Create orders partition
        partition_name := 'orders_' || year_val;
        EXECUTE format(
            'CREATE TABLE IF NOT EXISTS %I PARTITION OF orders FOR VALUES FROM (%L) TO (%L)',
            partition_name, start_date, end_date
        );
        RAISE NOTICE 'Created partition: %', partition_name;
    END LOOP;
END $$;

### 6. Maintenance Schedule (Recommended)

#### MONTHLY partitions (audit_logs, error_logs):
# Monthly (before month end):
-- Create partition for next month
-- Example: On Nov 25th, create December partition
CREATE TABLE audit_logs_2025_12 PARTITION OF audit_logs
    FOR VALUES FROM ('2025-12-01 00:00:00+00') TO ('2026-01-01 00:00:00+00');

CREATE TABLE error_logs_2025_12 PARTITION OF error_logs
    FOR VALUES FROM ('2025-12-01 00:00:00+00') TO ('2026-01-01 00:00:00+00');

#### YEARLY partitions (orders):
# Yearly (before year end):
-- Create partition for next year
-- Example: In November/December 2025, create 2026 partition
CREATE TABLE orders_2026 PARTITION OF orders
    FOR VALUES FROM ('2026-01-01 00:00:00+00') TO ('2027-01-01 00:00:00+00');

### 7. Quick Reference SQL File
# Create a file: scripts/manage-partitions.sql
-- ==============================================
-- PARTITION MANAGEMENT QUICK REFERENCE
-- ==============================================

-- 1. LIST ALL PARTITIONS
SELECT tablename, pg_size_pretty(pg_total_relation_size('public.'||tablename)) AS size
FROM pg_tables
WHERE tablename LIKE 'audit_logs_%' OR tablename LIKE 'error_logs_%' OR tablename LIKE 'orders_%'
ORDER BY tablename;

-- 2. CREATE NEXT MONTH PARTITION (Update dates as needed)
-- Example: Creating for December 2025 (MONTHLY: audit_logs, error_logs)
CREATE TABLE audit_logs_2025_12 PARTITION OF audit_logs
    FOR VALUES FROM ('2025-12-01 00:00:00+00') TO ('2026-01-01 00:00:00+00');

CREATE TABLE error_logs_2025_12 PARTITION OF error_logs
    FOR VALUES FROM ('2025-12-01 00:00:00+00') TO ('2026-01-01 00:00:00+00');

-- 3. CREATE NEXT YEAR PARTITION (Update dates as needed)
-- Example: Creating for 2027 (YEARLY: orders)
CREATE TABLE orders_2027 PARTITION OF orders
    FOR VALUES FROM ('2027-01-01 00:00:00+00') TO ('2028-01-01 00:00:00+00');

-- 4. CHECK DATA DISTRIBUTION
SELECT tableoid::regclass, COUNT(*) FROM audit_logs GROUP BY tableoid ORDER BY tableoid;
SELECT tableoid::regclass, COUNT(*) FROM error_logs GROUP BY tableoid ORDER BY tableoid;
SELECT tableoid::regclass, COUNT(*) FROM orders GROUP BY tableoid ORDER BY tableoid;

-- 5. ARCHIVE OLD PARTITION (Update partition name as needed)
-- MONTHLY: ALTER TABLE audit_logs DETACH PARTITION audit_logs_2024_01 CONCURRENTLY;
-- MONTHLY: DROP TABLE IF EXISTS audit_logs_2024_01;
-- YEARLY: ALTER TABLE orders DETACH PARTITION orders_2022 CONCURRENTLY;
-- YEARLY: DROP TABLE IF EXISTS orders_2022;

# Run it with:
docker exec -i $(docker-compose ps -q postgres) psql -U postgres -d appdb < scripts/manage-partitions.sql

# 8. Troubleshooting
# If partition creation fails:
-- Check if partition already exists
SELECT tablename FROM pg_tables WHERE tablename = 'audit_logs_2026_01' OR tablename = 'orders_2027';

-- Check current partition boundaries
SELECT pg_get_expr(relpartbound, oid) 
FROM pg_class 
WHERE relname LIKE 'audit_logs_%' OR relname LIKE 'error_logs_%' OR relname LIKE 'orders_%';

# If data goes to default partition:
-- Check default partition
SELECT COUNT(*) FROM audit_logs_default;
SELECT COUNT(*) FROM error_logs_default;
SELECT COUNT(*) FROM orders_default;

-- If it has data, you're missing partitions for those date ranges
-- Create the missing partitions

# 9. Partition Strategy Summary
-- audit_logs: MONTHLY partitions (high volume, frequent access)
-- error_logs: MONTHLY partitions (high volume, frequent access)
-- orders: YEARLY partitions (business data, long-term retention)
