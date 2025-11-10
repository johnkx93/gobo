-- name: GetPermissionsByRole :many
SELECT p.*
FROM permissions p
INNER JOIN role_permissions rp ON rp.permission_id = p.id
WHERE rp.role = $1 AND p.is_active = true;

-- name: GetAllPermissions :many
SELECT * FROM permissions
WHERE is_active = true
ORDER BY category, code;

-- name: GetPermissionByCode :one
SELECT * FROM permissions
WHERE code = $1 AND is_active = true;

-- name: CreatePermission :one
INSERT INTO permissions (code, name, description, category)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdatePermission :one
UPDATE permissions
SET name = COALESCE(sqlc.narg('name'), name),
    description = COALESCE(sqlc.narg('description'), description),
    category = COALESCE(sqlc.narg('category'), category),
    is_active = COALESCE(sqlc.narg('is_active'), is_active)
WHERE id = $1
RETURNING *;

-- name: DeletePermission :exec
DELETE FROM permissions WHERE id = $1;

-- name: AssignPermissionToRole :exec
INSERT INTO role_permissions (role, permission_id)
VALUES ($1, $2)
ON CONFLICT (role, permission_id) DO NOTHING;

-- name: RevokePermissionFromRole :exec
DELETE FROM role_permissions
WHERE role = $1 AND permission_id = $2;

-- name: GetRolePermissionCodes :many
SELECT p.code
FROM permissions p
INNER JOIN role_permissions rp ON rp.permission_id = p.id
WHERE rp.role = $1 AND p.is_active = true;
