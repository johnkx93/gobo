-- Drop all tables in reverse order (respecting foreign key constraints)

-- Drop menu and permissions tables
DROP TABLE IF EXISTS menu_items CASCADE;
DROP TABLE IF EXISTS role_permissions CASCADE;
DROP TABLE IF EXISTS permissions CASCADE;

-- Drop log tables (partitioned)
DROP TABLE IF EXISTS error_logs CASCADE;
DROP TABLE IF EXISTS audit_logs CASCADE;

-- Drop orders table (partitioned)
DROP TABLE IF EXISTS orders CASCADE;

-- Drop admin and user tables
DROP TABLE IF EXISTS admins CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- Drop enums
DROP TYPE IF EXISTS audit_action CASCADE;

-- Drop functions
DROP FUNCTION IF EXISTS update_updated_at_column() CASCADE;

-- Drop extensions (optional - be careful in shared databases)
-- DROP EXTENSION IF EXISTS pgcrypto CASCADE;
