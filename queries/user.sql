-- name: CreateUser :one
INSERT INTO users (email, name, password_hash)
VALUES ($1, $2, $3)
RETURNING id, created_at;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1; 

-- name: GetUserByEmail :one
SELECT id, email, name, password_hash, email_verified, created_at, updated_at
FROM users
WHERE email = $1 LIMIT 1;

-- name: GetUserById :one
SELECT id, email, name, password_hash, email_verified, created_at, updated_at
FROM users
WHERE id = $1 LIMIT 1;

-- name: UpdateUserEmail :one
UPDATE users
SET email = $2, updated_at = NOW()
WHERE id = $1
RETURNING id, email, name, password_hash, email_verified, created_at, updated_at;

-- name: UpdateUserPassword :one
UPDATE users
SET password_hash = $2, updated_at = NOW()
WHERE id = $1
RETURNING id, email, name, password_hash, email_verified, created_at, updated_at;

-- name: UpdateUserName :one
UPDATE users
SET name = $2, updated_at = NOW()
WHERE id = $1
RETURNING id, email, name, password_hash, email_verified, created_at, updated_at;

-- name: UpdateUserFull :one
UPDATE users
SET email = $2, name = $3, password_hash = $4, email_verified = $5, updated_at = NOW()
WHERE id = $1
RETURNING id, email, name, password_hash, email_verified, created_at, updated_at;

