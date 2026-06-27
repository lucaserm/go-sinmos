-- +goose Up
CREATE TYPE permission_type AS ENUM ('LEAVE', 'STAY');
CREATE TABLE IF NOT EXISTS permissions (
    id UUID PRIMARY KEY,
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    type permission_type NOT NULL,
    description TEXT NOT NULL,
    requested_at TIMESTAMPTZ NOT NULL,
    scheduled_for TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS permissions;
DROP TYPE IF EXISTS permission_type;