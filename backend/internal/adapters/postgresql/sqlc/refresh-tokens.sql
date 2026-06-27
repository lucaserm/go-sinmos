-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, user_id, expires_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetRefreshTokenByUserId :one
SELECT * FROM refresh_tokens WHERE user_id = $1 AND expires_at >= $2;

-- name: GetRefreshTokenByToken :one
SELECT * FROM refresh_tokens WHERE token = $1 AND expires_at >= $2;

-- name: DeleteRefreshTokenByUserId :exec
DELETE FROM refresh_tokens WHERE user_id = $1;
