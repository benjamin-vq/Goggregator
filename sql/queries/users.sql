-- name: CreateUser :one
INSERT INTO users (id, createdAt, updatedAt, name, apiKey)
VALUES ($1, $2, $3, $4, encode(sha256(random()::text::bytea), 'hex'))
RETURNING *;

-- name: GetUserByApiKey :one
SELECT * FROM users
WHERE apikey = $1;