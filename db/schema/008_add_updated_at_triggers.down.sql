-- Drop triggers
DROP TRIGGER IF EXISTS trigger_update_users_updated_at ON users;
DROP TRIGGER IF EXISTS trigger_update_orders_updated_at ON orders;
DROP TRIGGER IF EXISTS trigger_update_admins_updated_at ON admins;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();
