-- Convert orders table to partitioned table

-- Step 1: Drop existing indexes (they'll be recreated on the partitioned table)
DROP INDEX IF EXISTS idx_orders_user_id;
DROP INDEX IF EXISTS idx_orders_order_number;
DROP INDEX IF EXISTS idx_orders_status;
DROP INDEX IF EXISTS idx_orders_created_at;

-- Step 2: Rename existing orders table
ALTER TABLE orders RENAME TO orders_old;

-- Step 3: Create new partitioned orders table
CREATE TABLE orders (
    id UUID DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    order_number VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    total_amount DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    PRIMARY KEY (id, created_at)
) PARTITION BY RANGE (created_at);

-- Step 3: Create partitions for 4 months (November 2025 to February 2026)

-- November 2025
CREATE TABLE orders_2025_11 PARTITION OF orders
    FOR VALUES FROM ('2025-11-01 00:00:00+00') TO ('2025-12-01 00:00:00+00');

-- December 2025
CREATE TABLE orders_2025_12 PARTITION OF orders
    FOR VALUES FROM ('2025-12-01 00:00:00+00') TO ('2026-01-01 00:00:00+00');

-- January 2026
CREATE TABLE orders_2026_01 PARTITION OF orders
    FOR VALUES FROM ('2026-01-01 00:00:00+00') TO ('2026-02-01 00:00:00+00');

-- February 2026
CREATE TABLE orders_2026_02 PARTITION OF orders
    FOR VALUES FROM ('2026-02-01 00:00:00+00') TO ('2026-03-01 00:00:00+00');

-- Default partition for out-of-range dates
CREATE TABLE orders_default PARTITION OF orders DEFAULT;

-- Step 5: Create indexes on parent table (automatically propagates to all partitions)
CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_order_number ON orders(order_number);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_created_at ON orders(created_at DESC);

-- Step 6: Migrate data from old table to new partitioned table
INSERT INTO orders 
SELECT id, user_id, order_number, status, total_amount, notes, created_at, updated_at
FROM orders_old;

-- Step 6: Verify data migration
DO $$
DECLARE
    old_count INTEGER;
    new_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO old_count FROM orders_old;
    SELECT COUNT(*) INTO new_count FROM orders;
    
    IF old_count != new_count THEN
        RAISE EXCEPTION 'Data migration failed: old table has % rows, new table has % rows', 
            old_count, new_count;
    END IF;
    
    RAISE NOTICE 'Data migration successful: % rows migrated', new_count;
END $$;

-- Step 7: Drop old table
DROP TABLE orders_old;
