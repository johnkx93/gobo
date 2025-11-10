-- Remove default_address_id column from users table
DROP INDEX IF EXISTS idx_users_default_address_id;
ALTER TABLE users DROP COLUMN IF EXISTS default_address_id;
