-- name: CreateOccurrenceType :one
INSERT INTO occurrence_types (id, code, description, severity)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateOccurrenceType :one
UPDATE occurrence_types
SET code = $2, description = $3, severity = $4, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteOccurrenceType :exec
DELETE FROM occurrence_types WHERE id = $1;

-- name: GetAllOccurrenceTypes :many
SELECT * FROM occurrence_types LIMIT $1 OFFSET $2;

-- name: GetOccurrenceTypeByID :one
SELECT * FROM occurrence_types WHERE id = $1;