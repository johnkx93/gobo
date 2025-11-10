-- name: CreateAddress :one
INSERT INTO addresses (
    user_id,
    address,
    floor,
    unit_no,
    block_tower,
    company_name
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetAddressByID :one
SELECT * FROM addresses
WHERE id = $1 LIMIT 1;

-- name: GetAddressByIDAndUserID :one
SELECT * FROM addresses
WHERE id = $1 AND user_id = $2 LIMIT 1;

-- name: GetAddressesByUserID :many
SELECT * FROM addresses
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: ListAllAddresses :many
SELECT * FROM addresses
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountAddresses :one
SELECT COUNT(*) FROM addresses;

-- name: CountAddressesByUserID :one
SELECT COUNT(*) FROM addresses
WHERE user_id = $1;

-- name: UpdateAddress :one
UPDATE addresses
SET
    address = COALESCE(sqlc.narg('address'), address),
    floor = COALESCE(sqlc.narg('floor'), floor),
    unit_no = COALESCE(sqlc.narg('unit_no'), unit_no),
    block_tower = COALESCE(sqlc.narg('block_tower'), block_tower),
    company_name = COALESCE(sqlc.narg('company_name'), company_name),
    updated_at = CURRENT_TIMESTAMP
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: UpdateAddressForUser :one
UPDATE addresses
SET
    address = COALESCE(sqlc.narg('address'), address),
    floor = COALESCE(sqlc.narg('floor'), floor),
    unit_no = COALESCE(sqlc.narg('unit_no'), unit_no),
    block_tower = COALESCE(sqlc.narg('block_tower'), block_tower),
    company_name = COALESCE(sqlc.narg('company_name'), company_name),
    updated_at = CURRENT_TIMESTAMP
WHERE id = sqlc.arg('id') AND user_id = sqlc.arg('user_id')
RETURNING *;

-- name: DeleteAddress :exec
DELETE FROM addresses
WHERE id = $1;

-- name: DeleteAddressForUser :exec
DELETE FROM addresses
WHERE id = $1 AND user_id = $2;

-- name: DeleteAddressesByUserID :exec
DELETE FROM addresses
WHERE user_id = $1;

-- name: SetDefaultAddress :one
UPDATE users
SET default_address_id = $2
WHERE id = $1
RETURNING *;

-- name: SetDefaultAddressForUser :one
UPDATE users u
SET default_address_id = $2
WHERE u.id = $1
AND EXISTS (
    SELECT 1 FROM addresses a 
    WHERE a.id = $2 AND a.user_id = $1
)
RETURNING u.*;

-- name: ClearDefaultAddress :one
UPDATE users
SET default_address_id = NULL
WHERE id = $1
RETURNING *;

-- name: GetUserWithDefaultAddress :one
SELECT 
    u.*,
    a.id as default_address_id,
    a.address as default_address,
    a.floor as default_floor,
    a.unit_no as default_unit_no,
    a.block_tower as default_block_tower,
    a.company_name as default_company_name
FROM users u
LEFT JOIN addresses a ON u.default_address_id = a.id
WHERE u.id = $1
LIMIT 1;
