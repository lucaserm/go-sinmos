-- +goose Up
CREATE TABLE IF NOT EXISTS warnings (
    id UUID PRIMARY KEY,
    occurrence_id UUID UNIQUE NOT NULL REFERENCES occurrences(id),
    report TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- +goose Down
DROP TABLE IF EXISTS warnings;
