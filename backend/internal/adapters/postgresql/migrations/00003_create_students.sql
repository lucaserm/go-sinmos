-- +goose Up
CREATE TABLE IF NOT EXISTS students (
    id UUID PRIMARY KEY,
    cpf VARCHAR(30) UNIQUE NOT NULL,
    ra VARCHAR(20) UNIQUE NOT NULL,
    photo_url TEXT,
    name TEXT NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS students;
