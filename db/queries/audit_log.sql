-- name: CreateAuditLog :one
INSERT INTO audit_logs (
    user_id,
    action,
    entity_type,
    entity_id,
    old_data,
    new_data,
    request_id,
    ip_address,
    user_agent,
    metadata,
    created_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
) RETURNING *;

-- name: GetAuditLogByID :one
SELECT * FROM audit_logs
WHERE id = $1 LIMIT 1;

-- name: ListAuditLogsByEntity :many
SELECT * FROM audit_logs
WHERE entity_type = $1 AND entity_id = $2
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;

-- name: ListAuditLogsByUser :many
SELECT * FROM audit_logs
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListAuditLogsByEntityType :many
SELECT * FROM audit_logs
WHERE entity_type = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListAuditLogsByDateRange :many
SELECT * FROM audit_logs
WHERE created_at >= $1 AND created_at < $2
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;

-- name: ListAuditLogsByRequestID :many
SELECT * FROM audit_logs
WHERE request_id = $1
ORDER BY created_at ASC;

-- name: ListAuditLogsByEntityAndDateRange :many
SELECT * FROM audit_logs
WHERE entity_type = $1 
  AND entity_id = $2
  AND created_at >= $3 
  AND created_at < $4
ORDER BY created_at DESC
LIMIT $5 OFFSET $6;

-- name: CountAuditLogsByEntity :one
SELECT COUNT(*) FROM audit_logs
WHERE entity_type = $1 AND entity_id = $2;

-- name: CountAuditLogsByUser :one
SELECT COUNT(*) FROM audit_logs
WHERE user_id = $1;
