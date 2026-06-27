-- name: CreatePermission :one
INSERT INTO permissions (id, student_id, type, description, requested_at, scheduled_for)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdatePermission :one
UPDATE permissions
SET student_id = $2, type = $3, description = $4, requested_at = $5, scheduled_for = $6, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeletePermission :exec
DELETE FROM permissions WHERE id = $1;

-- name: GetAllPermissions :many
SELECT * FROM permissions LIMIT $1 OFFSET $2;

-- name: GetPermissionByID :one
SELECT * FROM permissions WHERE id = $1;