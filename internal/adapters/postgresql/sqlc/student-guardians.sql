-- name: AddStudentGuardian :exec
INSERT INTO students_guardians (student_id, guardian_id)
VALUES ($1, $2)
ON CONFLICT (student_id, guardian_id) DO NOTHING;

-- name: RemoveStudentGuardian :exec
DELETE FROM students_guardians
WHERE student_id = $1 AND guardian_id = $2;