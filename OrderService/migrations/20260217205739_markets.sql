-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS  markets(
id SERIAL PRIMARY KEY,
name TEXT,
market_id TEXT NOT NULL UNIQUE,
created_at TIMESTAMPTZ NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS markets;
-- +goose StatementEnd

