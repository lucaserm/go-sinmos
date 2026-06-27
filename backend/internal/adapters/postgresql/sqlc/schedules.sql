-- name: CreateSchedule :one
INSERT INTO schedules (id, subject_id, session, day_of_week, start_time, end_time)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetAllSchedules :many
SELECT * FROM schedules LIMIT $1 OFFSET $2;

-- name: GetScheduleByID :one
SELECT * FROM schedules WHERE id = $1;

-- name: UpdateSchedule :one
UPDATE schedules
SET subject_id = $2, session = $3, day_of_week = $4, start_time = $5, end_time = $6, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteSchedule :exec
DELETE FROM schedules WHERE id = $1;