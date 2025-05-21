CREATE TABLE users(
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    name TEXT NOT NULL,
    surname TEXT NOT NULL,
    patronymic TEXT
);

-- using migrate 
-- migrate -path sql/schema -database "DB_URL?sslmode=disable" up
-- example: 
-- migrate -path sql/schema -database "postgres://postgres:1111@localhost:5436/postgres?sslmode=disable" up
