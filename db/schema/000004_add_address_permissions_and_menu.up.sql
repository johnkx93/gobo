-- ==============================================
-- ADD ADDRESS MANAGEMENT PERMISSIONS
-- ==============================================

-- Insert address permissions
INSERT INTO permissions (code, name, description, category) VALUES
    ('addresses.create', 'Create Addresses', 'Ability to create new addresses for users', 'addresses'),
    ('addresses.read', 'Read Addresses', 'Ability to view address information', 'addresses'),
    ('addresses.update', 'Update Addresses', 'Ability to update address information', 'addresses'),
    ('addresses.delete', 'Delete Addresses', 'Ability to delete addresses', 'addresses');

-- ==============================================
-- ASSIGN PERMISSIONS TO SUPER ADMIN
-- ==============================================

-- Super Admin gets all address permissions
INSERT INTO role_permissions (role, permission_id)
SELECT 'super_admin', id FROM permissions 
WHERE code IN (
    'addresses.create',
    'addresses.read', 
    'addresses.update',
    'addresses.delete'
) AND is_active = true;

-- Admin gets read and update permissions
INSERT INTO role_permissions (role, permission_id)
SELECT 'admin', id FROM permissions 
WHERE code IN (
    'addresses.read',
    'addresses.update'
) AND is_active = true;

-- Moderator gets read-only permission
INSERT INTO role_permissions (role, permission_id)
SELECT 'moderator', id FROM permissions 
WHERE code = 'addresses.read' AND is_active = true;

-- ==============================================
-- ADD ADDRESS MENU ITEMS
-- ==============================================

-- Root menu item for Address Management
INSERT INTO menu_items (code, label, icon, order_index, permission_id) VALUES
    ('addresses', 'Address Management', 'map-pin', 6, NULL);

-- Child menu items for Address Management
INSERT INTO menu_items (parent_id, code, label, path, order_index, permission_id)
SELECT 
    (SELECT id FROM menu_items WHERE code = 'addresses'),
    'addresses-list',
    'Address List',
    '/admin/addresses',
    1,
    (SELECT id FROM permissions WHERE code = 'addresses.read')
UNION ALL
SELECT 
    (SELECT id FROM menu_items WHERE code = 'addresses'),
    'addresses-create',
    'Create Address',
    '/admin/addresses/create',
    2,
    (SELECT id FROM permissions WHERE code = 'addresses.create');
