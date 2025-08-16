-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS api_keys(
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL, 
  prefix TEXT NOT NULL,
  api_key TEXT UNIQUE NOT NULL,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  permissions TEXT[],
  metadata JSONB,
  expires_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ DEFAULT NOW()
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS api_keys;
-- +goose StatementEnd
