-- name: CreateGuardian :one
INSERT INTO guardians (id, phone, email, name) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetGuardians :many
SELECT * FROM guardians LIMIT $1 OFFSET $2;

-- name: GetGuardianByID :one
SELECT * FROM guardians WHERE id = $1;

-- name: UpdateGuardian :one
UPDATE guardians SET phone = $2, email = $3, name = $4, updated_at = now() WHERE id = $1 RETURNING *;

-- name: DeleteGuardian :exec
DELETE FROM guardians WHERE id = $1;