-- +goose Up
CREATE TABLE IF NOT EXISTS schedules (
    id UUID PRIMARY KEY,
    subject_id UUID NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
    session session NOT NULL,
    day_of_week VARCHAR(30) NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS schedules;
