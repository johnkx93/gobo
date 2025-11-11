-- Recreate Orders table, partitions, permissions and menu items

-- Recreate parent orders table (partitioned by created_at)
CREATE TABLE IF NOT EXISTS orders (
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

-- Recreate yearly partitions
CREATE TABLE IF NOT EXISTS orders_2025 PARTITION OF orders
    FOR VALUES FROM ('2025-01-01 00:00:00+00') TO ('2026-01-01 00:00:00+00');

CREATE TABLE IF NOT EXISTS orders_2026 PARTITION OF orders
    FOR VALUES FROM ('2026-01-01 00:00:00+00') TO ('2027-01-01 00:00:00+00');

CREATE TABLE IF NOT EXISTS orders_default PARTITION OF orders DEFAULT;

-- Indexes
CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);
CREATE INDEX IF NOT EXISTS idx_orders_order_number ON orders(order_number);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders(created_at DESC);

-- Recreate trigger to update updated_at
CREATE TRIGGER IF NOT EXISTS trigger_update_orders_updated_at
    BEFORE UPDATE ON orders
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Re-insert order-related permissions
INSERT INTO permissions (code, name, description, category)
VALUES
    ('orders.create', 'Create Orders', 'Ability to create new orders', 'orders'),
    ('orders.read', 'Read Orders', 'Ability to view order information', 'orders'),
    ('orders.update', 'Update Orders', 'Ability to update order information', 'orders'),
    ('orders.delete', 'Delete Orders', 'Ability to delete orders', 'orders')
ON CONFLICT (code) DO NOTHING;

-- Re-insert role_permissions mapping for order permissions
INSERT INTO role_permissions (role, permission_id)
SELECT 'super_admin', id FROM permissions WHERE category = 'orders'
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role, permission_id)
SELECT 'admin', id FROM permissions WHERE code IN ('orders.read', 'orders.update')
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role, permission_id)
SELECT 'moderator', id FROM permissions WHERE code = 'orders.read'
ON CONFLICT DO NOTHING;

-- Recreate menu items for orders (root and children)
INSERT INTO menu_items (code, label, icon, order_index, permission_id)
VALUES ('orders', 'Order Management', 'shopping-cart', 2, NULL)
ON CONFLICT (code) DO NOTHING;

-- Child menu items
WITH root AS (SELECT id FROM menu_items WHERE code = 'orders')
INSERT INTO menu_items (parent_id, code, label, path, order_index, permission_id)
SELECT root.id, 'orders-list', 'Order List', '/admin/orders', 1, (SELECT id FROM permissions WHERE code = 'orders.read')
FROM root
ON CONFLICT (code) DO NOTHING;

WITH root AS (SELECT id FROM menu_items WHERE code = 'orders')
INSERT INTO menu_items (parent_id, code, label, path, order_index, permission_id)
SELECT root.id, 'orders-update', 'Update Orders', '/admin/orders/bulk-update', 2, (SELECT id FROM permissions WHERE code = 'orders.update')
FROM root
ON CONFLICT (code) DO NOTHING;
