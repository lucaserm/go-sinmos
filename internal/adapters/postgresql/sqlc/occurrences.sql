-- name: CreateOccurrence :one
INSERT INTO occurrences (id, user_id, occurrence_type_id, student_id, occurred_at, user_related_id, status)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UpdateOccurrence :one
UPDATE occurrences
SET user_id = $2, occurrence_type_id = $3, student_id = $4, occurred_at = $5, user_related_id = $6, status = $7
WHERE id = $1
RETURNING *;

-- name: DeleteOccurrence :exec
DELETE FROM occurrences WHERE id = $1;

-- name: GetAllOccurrences :many
SELECT * FROM occurrences LIMIT $1 OFFSET $2;

-- name: GetOccurrenceByID :one
SELECT * FROM occurrences WHERE id = $1;