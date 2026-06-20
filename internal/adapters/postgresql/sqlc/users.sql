-- name: CreateUser :one
INSERT INTO users (id, name, code, hashed_password, role)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByCode :one
SELECT * FROM users WHERE code = $1;
