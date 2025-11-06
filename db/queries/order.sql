-- name: CreateOrder :one
INSERT INTO orders (
    user_id,
    order_number,
    status,
    total_amount,
    notes
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetOrderByID :one
SELECT * FROM orders
WHERE id = $1 LIMIT 1;

-- name: GetOrderByOrderNumber :one
SELECT * FROM orders
WHERE order_number = $1 LIMIT 1;

-- name: ListOrdersByUserID :many
SELECT * FROM orders
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListOrders :many
SELECT * FROM orders
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateOrderStatus :one
UPDATE orders
SET
    status = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: UpdateOrder :one
UPDATE orders
SET
    status = COALESCE(sqlc.narg('status'), status),
    total_amount = COALESCE(sqlc.narg('total_amount'), total_amount),
    notes = COALESCE(sqlc.narg('notes'), notes),
    updated_at = CURRENT_TIMESTAMP
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteOrder :exec
DELETE FROM orders
WHERE id = $1;
