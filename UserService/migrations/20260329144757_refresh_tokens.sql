-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS refresh_tokens(
    id SERIAL PRIMARY KEY,
    token TEXT NOT NULL,
    expire_at TIMESTAMPTZ NOT NULL,
    is_revoked BOOLEAN NOT NULL,
    user_id UUID REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL,
    update_at TIMESTAMPTZ 
);

CREATE INDEX idx_refresh_tokens_token ON refresh_tokens(token);
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS refresh_tokens;
-- +goose StatementEnd
