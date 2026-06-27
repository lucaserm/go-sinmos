-- name: CreateSubject :one
INSERT INTO subjects (id, course_id, name, semester, section)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetAllSubjects :many
SELECT * FROM subjects LIMIT $1 OFFSET $2;

-- name: GetSubjectByID :one
SELECT * FROM subjects WHERE id = $1;

-- name: UpdateSubject :one
UPDATE subjects
SET name = $2,
    semester = $3,
    section = $4,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteSubject :exec
DELETE FROM subjects WHERE id = $1;