-- Remove address menu items
DELETE FROM menu_items WHERE code IN ('addresses', 'addresses-list', 'addresses-create');

-- Remove address role permissions
DELETE FROM role_permissions 
WHERE permission_id IN (
    SELECT id FROM permissions WHERE category = 'addresses'
);

-- Remove address permissions
DELETE FROM permissions WHERE category = 'addresses';
