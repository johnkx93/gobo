-- Drop triggers
DROP TRIGGER IF EXISTS trigger_update_menu_items_updated_at ON menu_items;
DROP TRIGGER IF EXISTS trigger_update_permissions_updated_at ON permissions;

-- Drop indexes
DROP INDEX IF EXISTS idx_permissions_is_active;
DROP INDEX IF EXISTS idx_permissions_code;
DROP INDEX IF EXISTS idx_menu_items_is_active;
DROP INDEX IF EXISTS idx_menu_items_permission_id;
DROP INDEX IF EXISTS idx_menu_items_parent_id;
DROP INDEX IF EXISTS idx_role_permissions_permission_id;
DROP INDEX IF EXISTS idx_role_permissions_role;

-- Drop tables (order matters due to foreign keys)
DROP TABLE IF EXISTS menu_items;
DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS permissions;
