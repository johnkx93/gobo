-- name: GetMenuItemsByRole :many
WITH RECURSIVE menu_tree AS (
    -- Get root menu items (no parent)
    SELECT 
        mi.id,
        mi.parent_id,
        mi.code,
        mi.label,
        mi.icon,
        mi.path,
        mi.permission_id,
        mi.order_index,
        0 as depth
    FROM menu_items mi
    WHERE mi.parent_id IS NULL 
      AND mi.is_active = true
      AND (
          mi.permission_id IS NULL 
          OR mi.permission_id IN (
              SELECT rp.permission_id 
              FROM role_permissions rp 
              WHERE rp.role = $1
          )
      )
    
    UNION ALL
    
    -- Get child menu items recursively
    SELECT 
        mi.id,
        mi.parent_id,
        mi.code,
        mi.label,
        mi.icon,
        mi.path,
        mi.permission_id,
        mi.order_index,
        mt.depth + 1
    FROM menu_items mi
    INNER JOIN menu_tree mt ON mi.parent_id = mt.id
    WHERE mi.is_active = true
      AND (
          mi.permission_id IS NULL 
          OR mi.permission_id IN (
              SELECT rp.permission_id 
              FROM role_permissions rp 
              WHERE rp.role = $1
          )
      )
)
SELECT * FROM menu_tree
ORDER BY depth, order_index;

-- name: GetAllMenuItems :many
SELECT * FROM menu_items
WHERE is_active = true
ORDER BY order_index;

-- name: GetMenuItemByID :one
SELECT * FROM menu_items
WHERE id = $1 AND is_active = true;

-- name: GetMenuItemByCode :one
SELECT * FROM menu_items
WHERE code = $1 AND is_active = true;

-- name: CreateMenuItem :one
INSERT INTO menu_items (parent_id, code, label, icon, path, permission_id, order_index)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UpdateMenuItem :one
UPDATE menu_items
SET label = COALESCE(sqlc.narg('label'), label),
    icon = COALESCE(sqlc.narg('icon'), icon),
    path = COALESCE(sqlc.narg('path'), path),
    permission_id = COALESCE(sqlc.narg('permission_id'), permission_id),
    order_index = COALESCE(sqlc.narg('order_index'), order_index),
    is_active = COALESCE(sqlc.narg('is_active'), is_active)
WHERE id = $1
RETURNING *;

-- name: DeleteMenuItem :exec
DELETE FROM menu_items WHERE id = $1;

-- name: GetChildMenuItems :many
SELECT * FROM menu_items
WHERE parent_id = $1 AND is_active = true
ORDER BY order_index;
