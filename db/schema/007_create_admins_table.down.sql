-- Drop trigger and function
DROP TRIGGER IF EXISTS trigger_update_admins_updated_at ON admins;
DROP FUNCTION IF EXISTS update_admins_updated_at();

-- Drop indexes
DROP INDEX IF EXISTS idx_admins_role;
DROP INDEX IF EXISTS idx_admins_username;
DROP INDEX IF EXISTS idx_admins_email;

-- Drop admins table
DROP TABLE IF EXISTS admins;
