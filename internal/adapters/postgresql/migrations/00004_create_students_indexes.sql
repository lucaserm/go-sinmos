-- +goose Up
CREATE INDEX idx_students_cpf
ON students(cpf);

CREATE INDEX idx_students_ra
ON students(ra);

-- +goose Down
DROP INDEX IF EXISTS idx_students_cpf;
DROP INDEX IF EXISTS idx_students_ra;