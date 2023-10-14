-- +goose Up
CREATE TABLE IF NOT EXISTS chat (
    id SERIAL PRIMARY KEY,
    user_from varchar(255) not null,
    user_to varchar(255) not null,
    text text,
    timestamp timestamptz not null default now()
);

-- +goose Down
DROP TABLE IF EXISTS chat;
