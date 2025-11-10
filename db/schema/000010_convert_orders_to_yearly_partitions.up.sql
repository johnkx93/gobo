-- Convert orders table from monthly partitions to yearly partitions

-- Step 1: Detach the default partition from the current orders table
ALTER TABLE orders DETACH PARTITION orders_default;

-- Step 2: Rename current partitioned table
ALTER TABLE orders RENAME TO orders_monthly;

-- Step 3: Rename detached default partition to avoid conflict
ALTER TABLE orders_default RENAME TO orders_default_old;

-- Step 4: Create new partitioned orders table with yearly partitions
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

-- Step 5: Create yearly partitions for 2025 and 2026
CREATE TABLE orders_2025 PARTITION OF orders
    FOR VALUES FROM ('2025-01-01 00:00:00+00') TO ('2026-01-01 00:00:00+00');

CREATE TABLE orders_2026 PARTITION OF orders
    FOR VALUES FROM ('2026-01-01 00:00:00+00') TO ('2027-01-01 00:00:00+00');

-- Default partition for out-of-range dates
CREATE TABLE orders_default PARTITION OF orders DEFAULT;

-- Step 6: Create indexes on parent table (automatically propagates to all partitions)
-- Use IF NOT EXISTS to avoid conflicts if indexes were inherited
CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);
CREATE INDEX IF NOT EXISTS idx_orders_order_number ON orders(order_number);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders(created_at DESC);

-- Step 7: Migrate data from monthly partitions to yearly partitions
INSERT INTO orders 
SELECT id, user_id, order_number, status, total_amount, notes, created_at, updated_at
FROM orders_monthly;

-- Step 8: Migrate any data from old default partition if it exists
INSERT INTO orders 
SELECT id, user_id, order_number, status, total_amount, notes, created_at, updated_at
FROM orders_default_old
WHERE EXISTS (SELECT 1 FROM orders_default_old LIMIT 1);

-- Step 9: Verify data migration
DO $$
DECLARE
    old_count INTEGER;
    old_default_count INTEGER;
    new_count INTEGER;
    expected_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO old_count FROM orders_monthly;
    SELECT COUNT(*) INTO old_default_count FROM orders_default_old;
    SELECT COUNT(*) INTO new_count FROM orders;
    expected_count := old_count + old_default_count;
    
    IF expected_count != new_count THEN
        RAISE EXCEPTION 'Data migration failed: monthly table has % rows, old default has % rows, new table has % rows (expected %)', 
            old_count, old_default_count, new_count, expected_count;
    END IF;
    
    RAISE NOTICE 'Data migration successful: % rows migrated from monthly to yearly partitions (% from partitions + % from default)', 
        new_count, old_count, old_default_count;
END $$;

-- Step 10: Drop old monthly partitioned table (this will drop all monthly partitions)
DROP TABLE orders_monthly CASCADE;

-- Step 11: Drop old default partition
DROP TABLE IF EXISTS orders_default_old;

-- Step 12: Add trigger for auto-updating updated_at
CREATE TRIGGER trigger_update_orders_updated_at
    BEFORE UPDATE ON orders
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
