-- name: CreateAdmin :one
INSERT INTO admins (email, username, password_hash, first_name, last_name, role, is_active)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetAdminByID :one
SELECT * FROM admins
WHERE id = $1 AND is_active = true
LIMIT 1;

-- name: GetAdminByEmail :one
SELECT * FROM admins
WHERE email = $1 AND is_active = true
LIMIT 1;

-- name: GetAdminByUsername :one
SELECT * FROM admins
WHERE username = $1 AND is_active = true
LIMIT 1;

-- name: ListAdmins :many
SELECT * FROM admins
WHERE is_active = true
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateAdmin :one
UPDATE admins
SET 
    email = COALESCE($2, email),
    username = COALESCE($3, username),
    password_hash = COALESCE($4, password_hash),
    first_name = $5,
    last_name = $6,
    role = COALESCE($7, role),
    is_active = COALESCE($8, is_active),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteAdmin :exec
UPDATE admins
SET is_active = false, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: HardDeleteAdmin :exec
DELETE FROM admins
WHERE id = $1;
