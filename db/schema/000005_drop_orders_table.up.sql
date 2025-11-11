-- Drop Orders table and related permissions/menu items

-- Remove role_permissions entries for order permissions
DELETE FROM role_permissions WHERE permission_id IN (SELECT id FROM permissions WHERE category = 'orders');

-- Remove child menu items and root menu item for orders
DELETE FROM menu_items WHERE code LIKE 'orders%';

-- Remove order-related permissions
DELETE FROM permissions WHERE category = 'orders';

-- Drop trigger on orders (if exists)
DROP TRIGGER IF EXISTS trigger_update_orders_updated_at ON orders;

-- Drop partitions (if they exist)
DROP TABLE IF EXISTS orders_2025;
DROP TABLE IF EXISTS orders_2026;
DROP TABLE IF EXISTS orders_default;

-- Finally drop the parent orders table
DROP TABLE IF EXISTS orders;
