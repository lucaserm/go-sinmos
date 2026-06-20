-- +goose Up
CREATE TYPE session AS ENUM ('MORNING', 'AFTERNOON', 'EVENING');
CREATE TABLE IF NOT EXISTS courses (
    id UUID PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    description TEXT,
    session session NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS courses;
DROP TYPE IF EXISTS session;
