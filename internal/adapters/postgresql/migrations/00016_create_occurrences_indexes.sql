-- +goose Up
CREATE INDEX idx_occurrences_student
ON occurrences(student_id);

CREATE INDEX idx_occurrences_status
ON occurrences(status);

-- +goose Down
DROP INDEX IF EXISTS idx_occurrences_student;
DROP INDEX IF EXISTS idx_occurrences_status;
