-- name: CreateAPIKey :one
INSERT INTO api_keys (name, prefix, api_key, user_id, permissions, expires_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, created_at;

-- name: ListAPIKeys :many
SELECT id, name, prefix, api_key, user_id, permissions, expires_at, created_at
FROM api_keys
WHERE user_id = $1
LIMIT $2 OFFSET $3;

-- name: DeleteAPIKey :execrows
DELETE FROM api_keys where id = $1 and user_id = $2;

-- name: CheckAPIKeyValid :one
SELECT EXISTS (
    SELECT 1 FROM api_keys WHERE api_key = $1
) AS valid;

-- name: GetAPIKey :one
SELECT * FROM api_keys WHERE api_key = $1;
