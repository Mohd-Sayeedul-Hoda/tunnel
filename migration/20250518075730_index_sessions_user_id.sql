-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_sessions_user_id;
-- +goose StatementEnd
