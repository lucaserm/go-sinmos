-- +goose Up
CREATE TYPE role AS ENUM ('ADMIN', 'SUPPORT', 'RECEPTION');
CREATE TABLE IF NOT EXISTS users (
	id UUID PRIMARY KEY,
    name varchar(50) NOT NULL,
    code varchar(20) UNIQUE NOT NULL,
    hashed_password varchar(255) NOT NULL,
    role role NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);


-- +goose Down
DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS role;
