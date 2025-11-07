-- Drop error_logs table (partitions are automatically dropped)
DROP TABLE IF EXISTS error_logs CASCADE;

-- Drop audit_logs table (partitions are automatically dropped)
DROP TABLE IF EXISTS audit_logs CASCADE;

-- Drop audit action enum
DROP TYPE IF EXISTS audit_action;
