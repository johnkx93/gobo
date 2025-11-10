-- Create permissions table
CREATE TABLE permissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(50) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Create role_permissions table (many-to-many)
CREATE TABLE role_permissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    role VARCHAR(50) NOT NULL,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE(role, permission_id)
);

-- Create menu_items table
CREATE TABLE menu_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
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

-- Create indexes
CREATE INDEX idx_role_permissions_role ON role_permissions(role);
CREATE INDEX idx_role_permissions_permission_id ON role_permissions(permission_id);
CREATE INDEX idx_menu_items_parent_id ON menu_items(parent_id);
CREATE INDEX idx_menu_items_permission_id ON menu_items(permission_id);
CREATE INDEX idx_menu_items_is_active ON menu_items(is_active);
CREATE INDEX idx_permissions_code ON permissions(code);
CREATE INDEX idx_permissions_is_active ON permissions(is_active);

-- Add triggers for updated_at
CREATE TRIGGER trigger_update_permissions_updated_at
    BEFORE UPDATE ON permissions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_update_menu_items_updated_at
    BEFORE UPDATE ON menu_items
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Insert default permissions
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

-- Assign permissions to roles
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

-- Insert menu items
-- First, insert root menu items
INSERT INTO menu_items (code, label, icon, order_index, permission_id) VALUES
    ('users', 'User Management', 'users', 1, NULL),
    ('orders', 'Order Management', 'shopping-cart', 2, NULL),
    ('admins', 'Admin Management', 'shield', 3, (SELECT id FROM permissions WHERE code = 'admins.manage')),
    ('settings', 'Settings', 'settings', 4, NULL),
    ('analytics', 'Analytics', 'chart-bar', 5, NULL);

-- Insert child menu items for User Management
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

-- Insert child menu items for Order Management
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

-- Insert child menu items for Admin Management (single item, no parent needed)
INSERT INTO menu_items (parent_id, code, label, path, order_index, permission_id)
SELECT 
    (SELECT id FROM menu_items WHERE code = 'admins'),
    'admins-list',
    'Admin List',
    '/admin/admins',
    1,
    (SELECT id FROM permissions WHERE code = 'admins.read');

-- Insert child menu items for Settings
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

-- Insert child menu items for Analytics
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
