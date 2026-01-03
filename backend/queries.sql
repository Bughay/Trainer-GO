-- name: CreateUser :one
INSERT INTO users (username, hashed_password)
VALUES ($1, $2)
RETURNING user_id, username, created_at;

-- name: GetUserByUsername :one
SELECT user_id, hashed_password
FROM users
WHERE username = $1;

