-- name: AddStudentSubject :exec
INSERT INTO students_subjects (student_id, subject_id)
VALUES ($1, $2)
ON CONFLICT (student_id, subject_id) DO NOTHING;

-- name: RemoveStudentSubject :exec
DELETE FROM students_subjects
WHERE student_id = $1 AND subject_id = $2;