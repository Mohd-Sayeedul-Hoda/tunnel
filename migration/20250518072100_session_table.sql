-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS sessions(
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  token TEXT NOT NULL,
  expires_at TIMESTAMP WITH TIME ZONE,
  ip_address VARCHAR(255) NOT NULL,
  user_agent VARCHAR(1024) NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS sessions;

DROP INDEX IF EXISTS idx_sessions_user_id;
-- +goose StatementEnd
