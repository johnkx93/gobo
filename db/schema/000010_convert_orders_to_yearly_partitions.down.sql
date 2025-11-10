-- Revert orders table back to monthly partitions

-- Step 1: Detach the default partition from the current orders table
ALTER TABLE orders DETACH PARTITION orders_default;

-- Step 2: Rename yearly partitioned table
ALTER TABLE orders RENAME TO orders_yearly;

-- Step 3: Rename detached default partition to avoid conflict
ALTER TABLE orders_default RENAME TO orders_default_old;

-- Step 4: Create new partitioned orders table with monthly partitions
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

-- Step 3: Recreate monthly partitions (November 2025 to February 2026)
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

-- Step 4: Recreate indexes
CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);
CREATE INDEX IF NOT EXISTS idx_orders_order_number ON orders(order_number);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders(created_at DESC);

-- Step 5: Migrate data back from yearly partitions to monthly partitions
INSERT INTO orders 
SELECT id, user_id, order_number, status, total_amount, notes, created_at, updated_at
FROM orders_yearly;

-- Step 6: Migrate any data from old default partition if it exists
INSERT INTO orders 
SELECT id, user_id, order_number, status, total_amount, notes, created_at, updated_at
FROM orders_default_old
WHERE EXISTS (SELECT 1 FROM orders_default_old LIMIT 1);

-- Step 7: Verify data migration
DO $$
DECLARE
    old_count INTEGER;
    old_default_count INTEGER;
    new_count INTEGER;
    expected_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO old_count FROM orders_yearly;
    SELECT COUNT(*) INTO old_default_count FROM orders_default_old;
    SELECT COUNT(*) INTO new_count FROM orders;
    expected_count := old_count + old_default_count;
    
    IF expected_count != new_count THEN
        RAISE EXCEPTION 'Rollback failed: yearly table has % rows, old default has % rows, monthly table has % rows (expected %)', 
            old_count, old_default_count, new_count, expected_count;
    END IF;
    
    RAISE NOTICE 'Rollback successful: % rows migrated back to monthly partitions (% from partitions + % from default)', 
        new_count, old_count, old_default_count;
END $$;

-- Step 8: Drop yearly partitioned table (this will drop all yearly partitions)
DROP TABLE orders_yearly CASCADE;

-- Step 9: Drop old default partition
DROP TABLE IF EXISTS orders_default_old;

-- Step 10: Add trigger for auto-updating updated_at
CREATE TRIGGER trigger_update_orders_updated_at
    BEFORE UPDATE ON orders
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
