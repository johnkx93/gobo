-- name: CreateErrorLog :one
INSERT INTO error_logs (
    user_id,
    request_id,
    error_type,
    error_message,
    stack_trace,
    request_path,
    request_method,
    ip_address,
    user_agent,
    metadata,
    created_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
) RETURNING *;

-- name: GetErrorLogByID :one
SELECT * FROM error_logs
WHERE id = $1 LIMIT 1;

-- name: ListErrorLogsByRequestID :many
SELECT * FROM error_logs
WHERE request_id = $1
ORDER BY created_at ASC;

-- name: ListErrorLogsByUser :many
SELECT * FROM error_logs
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListErrorLogsByType :many
SELECT * FROM error_logs
WHERE error_type = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListErrorLogsByDateRange :many
SELECT * FROM error_logs
WHERE created_at >= $1 AND created_at < $2
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;

-- name: ListErrorLogsByPath :many
SELECT * FROM error_logs
WHERE request_path = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListRecentErrors :many
SELECT * FROM error_logs
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountErrorLogsByType :one
SELECT COUNT(*) FROM error_logs
WHERE error_type = $1;

-- name: CountErrorLogsByDateRange :one
SELECT COUNT(*) FROM error_logs
WHERE created_at >= $1 AND created_at < $2;
