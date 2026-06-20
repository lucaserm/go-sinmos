-- name: CreateWarning :one
INSERT INTO warnings (id, occurrence_id, report) VALUES ($1, $2, $3) RETURNING *;

-- name: GetAllWarnings :many
SELECT * FROM warnings LIMIT $1 OFFSET $2;

-- name: GetWarningByID :one
SELECT * FROM warnings WHERE id = $1;

-- name: UpdateWarning :one
UPDATE warnings SET occurrence_id = $2, report = $3, updated_at = now() WHERE id = $1 RETURNING *;

-- name: DeleteWarning :exec
DELETE FROM warnings WHERE id = $1;