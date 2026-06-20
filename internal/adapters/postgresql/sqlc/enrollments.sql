-- name: CreateEnrollment :one
INSERT INTO enrollments (id, student_id, course_id, year) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetAllEnrollments :many
SELECT * FROM enrollments LIMIT $1 OFFSET $2;

-- name: GetEnrollmentByID :one
SELECT * FROM enrollments WHERE id = $1;

-- name: UpdateEnrollment :one
UPDATE enrollments SET student_id = $2, course_id = $3, year = $4, updated_at = now() WHERE id = $1 RETURNING *;

-- name: DeleteEnrollment :exec
DELETE FROM enrollments WHERE id = $1;