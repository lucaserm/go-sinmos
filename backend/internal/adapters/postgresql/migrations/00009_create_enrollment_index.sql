-- +goose Up
CREATE INDEX idx_enrollments_student
ON enrollments(student_id);

-- +goose Down
DROP INDEX idx_enrollments_student;