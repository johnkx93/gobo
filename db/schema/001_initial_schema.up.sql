-- ==============================================
-- INITIAL DATABASE SCHEMA
-- ==============================================

-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- ==============================================
-- REUSABLE FUNCTIONS
-- ==============================================

-- Function to automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- ==============================================
-- ENUMS
-- ==============================================

CREATE TYPE audit_action AS ENUM ('CREATE', 'UPDATE', 'DELETE');

-- ==============================================
-- USERS TABLE
-- ==============================================

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);

CREATE TRIGGER trigger_update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- ==============================================
-- ADMINS TABLE
-- ==============================================

CREATE TABLE admins (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    role VARCHAR(50) NOT NULL DEFAULT 'admin',
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX idx_admins_email ON admins(email);
CREATE INDEX idx_admins_username ON admins(username);
CREATE INDEX idx_admins_role ON admins(role);

CREATE TRIGGER trigger_update_admins_updated_at
    BEFORE UPDATE ON admins
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Insert default super admin (password: 'admin123' - CHANGE THIS IN PRODUCTION!)
INSERT INTO admins (email, username, password_hash, first_name, last_name, role)
VALUES (
    'admin@example.com',
    'superadmin',
    '$2a$10$WfaDveNY3z.lAO6pk9vdgunHpQEowC1C5jsCLMF6qCIkvrwmEXlMW',
    'Super',
    'Admin',
    'super_admin'
);

-- ==============================================
-- ORDERS TABLE (YEARLY PARTITIONS)
-- ==============================================

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

-- Create yearly partitions for 2025 and 2026
CREATE TABLE orders_2025 PARTITION OF orders
    FOR VALUES FROM ('2025-01-01 00:00:00+00') TO ('2026-01-01 00:00:00+00');

CREATE TABLE orders_2026 PARTITION OF orders
    FOR VALUES FROM ('2026-01-01 00:00:00+00') TO ('2027-01-01 00:00:00+00');

CREATE TABLE orders_default PARTITION OF orders DEFAULT;

-- Create indexes
CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_order_number ON orders(order_number);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_created_at ON orders(created_at DESC);

CREATE TRIGGER trigger_update_orders_updated_at
    BEFORE UPDATE ON orders
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- ==============================================
-- AUDIT LOGS TABLE (MONTHLY PARTITIONS)
-- ==============================================

CREATE TABLE audit_logs (
    id UUID DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action audit_action NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id UUID NOT NULL,
    old_data JSONB,
    new_data JSONB,
    request_id VARCHAR(100),
    ip_address VARCHAR(45),
    user_agent TEXT,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    PRIMARY KEY (id, created_at)
) PARTITION BY RANGE (created_at);

-- Create monthly partitions
CREATE TABLE audit_logs_2025_11 PARTITION OF audit_logs
    FOR VALUES FROM ('2025-11-01 00:00:00+00') TO ('2025-12-01 00:00:00+00');

CREATE TABLE audit_logs_2025_12 PARTITION OF audit_logs
    FOR VALUES FROM ('2025-12-01 00:00:00+00') TO ('2026-01-01 00:00:00+00');

CREATE TABLE audit_logs_2026_01 PARTITION OF audit_logs
    FOR VALUES FROM ('2026-01-01 00:00:00+00') TO ('2026-02-01 00:00:00+00');

CREATE TABLE audit_logs_2026_02 PARTITION OF audit_logs
    FOR VALUES FROM ('2026-02-01 00:00:00+00') TO ('2026-03-01 00:00:00+00');

CREATE TABLE audit_logs_default PARTITION OF audit_logs DEFAULT;

-- Create indexes
CREATE INDEX idx_audit_logs_entity ON audit_logs(entity_type, entity_id);
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at DESC);
CREATE INDEX idx_audit_logs_request_id ON audit_logs(request_id);
CREATE INDEX idx_audit_logs_composite ON audit_logs(entity_type, entity_id, created_at DESC);

-- ==============================================
-- ERROR LOGS TABLE (MONTHLY PARTITIONS)
-- ==============================================

CREATE TABLE error_logs (
    id UUID DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    request_id VARCHAR(100),
    error_type VARCHAR(100) NOT NULL,
    error_message TEXT NOT NULL,
    stack_trace TEXT,
    request_path VARCHAR(255),
    request_method VARCHAR(10),
    ip_address VARCHAR(45),
    user_agent TEXT,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    PRIMARY KEY (id, created_at)
) PARTITION BY RANGE (created_at);

-- Create monthly partitions
CREATE TABLE error_logs_2025_11 PARTITION OF error_logs
    FOR VALUES FROM ('2025-11-01 00:00:00+00') TO ('2025-12-01 00:00:00+00');

CREATE TABLE error_logs_2025_12 PARTITION OF error_logs
    FOR VALUES FROM ('2025-12-01 00:00:00+00') TO ('2026-01-01 00:00:00+00');

CREATE TABLE error_logs_2026_01 PARTITION OF error_logs
    FOR VALUES FROM ('2026-01-01 00:00:00+00') TO ('2026-02-01 00:00:00+00');

CREATE TABLE error_logs_2026_02 PARTITION OF error_logs
    FOR VALUES FROM ('2026-02-01 00:00:00+00') TO ('2026-03-01 00:00:00+00');

CREATE TABLE error_logs_default PARTITION OF error_logs DEFAULT;

-- Create indexes
CREATE INDEX idx_error_logs_user_id ON error_logs(user_id);
CREATE INDEX idx_error_logs_request_id ON error_logs(request_id);
CREATE INDEX idx_error_logs_created_at ON error_logs(created_at DESC);
CREATE INDEX idx_error_logs_error_type ON error_logs(error_type);
CREATE INDEX idx_error_logs_request_path ON error_logs(request_path);

-- ==============================================
-- PERMISSIONS & MENU SYSTEM
-- ==============================================

-- Permissions table
CREATE TABLE permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(50) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX idx_permissions_code ON permissions(code);
CREATE INDEX idx_permissions_is_active ON permissions(is_active);

CREATE TRIGGER trigger_update_permissions_updated_at
    BEFORE UPDATE ON permissions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Role permissions table (many-to-many)
CREATE TABLE role_permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    role VARCHAR(50) NOT NULL,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE(role, permission_id)
);

CREATE INDEX idx_role_permissions_role ON role_permissions(role);
CREATE INDEX idx_role_permissions_permission_id ON role_permissions(permission_id);

-- Menu items table
CREATE TABLE menu_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    parent_id UUID REFERENCES menu_items(id) ON DELETE CASCADE,
    code VARCHAR(100) NOT NULL UNIQUE,
    label VARCHAR(255) NOT NULL,
    icon VARCHAR(50),
    path VARCHAR(255),
    permission_id UUID REFERENCES permissions(id) ON DELETE SET NULL,
    order_index INT NOT NULL DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX idx_menu_items_parent_id ON menu_items(parent_id);
CREATE INDEX idx_menu_items_permission_id ON menu_items(permission_id);
CREATE INDEX idx_menu_items_is_active ON menu_items(is_active);

CREATE TRIGGER trigger_update_menu_items_updated_at
    BEFORE UPDATE ON menu_items
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- ==============================================
-- SEED DATA: PERMISSIONS
-- ==============================================

INSERT INTO permissions (code, name, description, category) VALUES
    -- User permissions
    ('users.create', 'Create Users', 'Ability to create new users', 'users'),
    ('users.read', 'Read Users', 'Ability to view user information', 'users'),
    ('users.update', 'Update Users', 'Ability to update user information', 'users'),
    ('users.delete', 'Delete Users', 'Ability to delete users', 'users'),
    
    -- Order permissions
    ('orders.create', 'Create Orders', 'Ability to create new orders', 'orders'),
    ('orders.read', 'Read Orders', 'Ability to view order information', 'orders'),
    ('orders.update', 'Update Orders', 'Ability to update order information', 'orders'),
    ('orders.delete', 'Delete Orders', 'Ability to delete orders', 'orders'),
    
    -- Admin permissions
    ('admins.create', 'Create Admins', 'Ability to create new admin accounts', 'admins'),
    ('admins.read', 'Read Admins', 'Ability to view admin information', 'admins'),
    ('admins.update', 'Update Admins', 'Ability to update admin information', 'admins'),
    ('admins.delete', 'Delete Admins', 'Ability to delete admin accounts', 'admins'),
    ('admins.manage', 'Manage Admins', 'Full admin management access', 'admins'),
    
    -- Settings permissions
    ('settings.general', 'General Settings', 'Ability to manage general settings', 'settings'),
    ('settings.security', 'Security Settings', 'Ability to manage security settings', 'settings'),
    
    -- Analytics permissions
    ('analytics.dashboard', 'Analytics Dashboard', 'Ability to view analytics dashboard', 'analytics'),
    ('analytics.reports', 'Analytics Reports', 'Ability to view and generate reports', 'analytics');

-- ==============================================
-- SEED DATA: ROLE PERMISSIONS
-- ==============================================

-- Super Admin - all permissions
INSERT INTO role_permissions (role, permission_id)
SELECT 'super_admin', id FROM permissions WHERE is_active = true;

-- Admin - most permissions except admin management and security
INSERT INTO role_permissions (role, permission_id)
SELECT 'admin', id FROM permissions 
WHERE code IN (
    'users.create', 'users.read', 'users.update',
    'orders.create', 'orders.read', 'orders.update',
    'settings.general',
    'analytics.dashboard'
) AND is_active = true;

-- Moderator - read-only permissions
INSERT INTO role_permissions (role, permission_id)
SELECT 'moderator', id FROM permissions 
WHERE code IN (
    'users.read',
    'orders.read',
    'analytics.dashboard'
) AND is_active = true;

-- ==============================================
-- SEED DATA: MENU ITEMS
-- ==============================================

-- Root menu items
INSERT INTO menu_items (code, label, icon, order_index, permission_id) VALUES
    ('users', 'User Management', 'users', 1, NULL),
    ('orders', 'Order Management', 'shopping-cart', 2, NULL),
    ('admins', 'Admin Management', 'shield', 3, (SELECT id FROM permissions WHERE code = 'admins.manage')),
    ('settings', 'Settings', 'settings', 4, NULL),
    ('analytics', 'Analytics', 'chart-bar', 5, NULL);

-- Child menu items for User Management
INSERT INTO menu_items (parent_id, code, label, path, order_index, permission_id)
SELECT 
    (SELECT id FROM menu_items WHERE code = 'users'),
    'users-create',
    'Create User',
    '/admin/users/create',
    1,
    (SELECT id FROM permissions WHERE code = 'users.create')
UNION ALL
SELECT 
    (SELECT id FROM menu_items WHERE code = 'users'),
    'users-list',
    'User List',
    '/admin/users',
    2,
    (SELECT id FROM permissions WHERE code = 'users.read');

-- Child menu items for Order Management
INSERT INTO menu_items (parent_id, code, label, path, order_index, permission_id)
SELECT 
    (SELECT id FROM menu_items WHERE code = 'orders'),
    'orders-list',
    'Order List',
    '/admin/orders',
    1,
    (SELECT id FROM permissions WHERE code = 'orders.read')
UNION ALL
SELECT 
    (SELECT id FROM menu_items WHERE code = 'orders'),
    'orders-update',
    'Update Orders',
    '/admin/orders/bulk-update',
    2,
    (SELECT id FROM permissions WHERE code = 'orders.update');

-- Child menu items for Admin Management
INSERT INTO menu_items (parent_id, code, label, path, order_index, permission_id)
SELECT 
    (SELECT id FROM menu_items WHERE code = 'admins'),
    'admins-list',
    'Admin List',
    '/admin/admins',
    1,
    (SELECT id FROM permissions WHERE code = 'admins.read');

-- Child menu items for Settings
INSERT INTO menu_items (parent_id, code, label, path, order_index, permission_id)
SELECT 
    (SELECT id FROM menu_items WHERE code = 'settings'),
    'settings-general',
    'General Settings',
    '/admin/settings/general',
    1,
    (SELECT id FROM permissions WHERE code = 'settings.general')
UNION ALL
SELECT 
    (SELECT id FROM menu_items WHERE code = 'settings'),
    'settings-security',
    'Security',
    '/admin/settings/security',
    2,
    (SELECT id FROM permissions WHERE code = 'settings.security');

-- Child menu items for Analytics
INSERT INTO menu_items (parent_id, code, label, path, order_index, permission_id)
SELECT 
    (SELECT id FROM menu_items WHERE code = 'analytics'),
    'analytics-dashboard',
    'Dashboard',
    '/admin/analytics/dashboard',
    1,
    (SELECT id FROM permissions WHERE code = 'analytics.dashboard')
UNION ALL
SELECT 
    (SELECT id FROM menu_items WHERE code = 'analytics'),
    'analytics-reports',
    'Reports',
    '/admin/analytics/reports',
    2,
    (SELECT id FROM permissions WHERE code = 'analytics.reports');
