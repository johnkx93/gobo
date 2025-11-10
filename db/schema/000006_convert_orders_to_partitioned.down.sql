-- Revert orders table back to non-partitioned

-- Step 1: Rename partitioned table
ALTER TABLE orders RENAME TO orders_partitioned;

-- Step 2: Recreate non-partitioned orders table
CREATE TABLE orders (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    order_number VARCHAR(50) NOT NULL UNIQUE,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    total_amount DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Step 3: Recreate indexes
CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_order_number ON orders(order_number);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_created_at ON orders(created_at DESC);

-- Step 4: Migrate data back from partitioned table
INSERT INTO orders 
SELECT id, user_id, order_number, status, total_amount, notes, created_at, updated_at
FROM orders_partitioned;

-- Step 5: Verify data migration
DO $$
DECLARE
    old_count INTEGER;
    new_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO old_count FROM orders_partitioned;
    SELECT COUNT(*) INTO new_count FROM orders;
    
    IF old_count != new_count THEN
        RAISE EXCEPTION 'Data migration failed: partitioned table has % rows, new table has % rows', 
            old_count, new_count;
    END IF;
    
    RAISE NOTICE 'Rollback successful: % rows migrated back', new_count;
END $$;

-- Step 6: Drop partitioned table (this will drop all partitions)
DROP TABLE orders_partitioned CASCADE;
