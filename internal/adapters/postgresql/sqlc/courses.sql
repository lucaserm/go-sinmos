-- name: CreateCourse :one
INSERT INTO courses (id, name, description, session)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetCourseByID :one
SELECT * FROM courses WHERE id = $1;

-- name: GetAllCourses :many
SELECT * FROM courses LIMIT $1 OFFSET $2;

-- name: UpdateCourse :one
UPDATE courses
SET name = $2, description = $3, session = $4, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteCourse :exec
DELETE FROM courses WHERE id = $1;