-- +goose Up
CREATE TABLE IF NOT EXISTS students_subjects (
    subject_id UUID NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (student_id, subject_id)
);

-- +goose Down
DROP TABLE IF EXISTS students_subjects;
