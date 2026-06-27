-- +goose Up
CREATE TYPE occurrence_status AS ENUM ('PENDING', 'APPROVED', 'REPROVED');
CREATE TABLE IF NOT EXISTS occurrences (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    occurrence_type_id UUID NOT NULL REFERENCES occurrence_types(id),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    occurred_at TIMESTAMPTZ NOT NULL,
    user_related_id UUID REFERENCES users(id),
    status occurrence_status NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);


-- +goose Down
DROP TABLE IF EXISTS occurrences;
DROP TYPE IF EXISTS occurrence_status;
