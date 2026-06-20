-- name: CreateStudent :one
INSERT INTO students (id, cpf, ra, photo_url, name, email) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: UpdateStudent :one
UPDATE students SET cpf = $2, ra = $3, photo_url = $4, name = $5, email = $6, updated_at = now() WHERE id = $1 RETURNING *;

-- name: DeleteStudent :exec
DELETE FROM students WHERE id = $1;

-- name: GetAllStudents :many
SELECT * FROM students LIMIT $1 OFFSET $2;

-- name: GetStudentByID :one
SELECT * FROM students WHERE id = $1;