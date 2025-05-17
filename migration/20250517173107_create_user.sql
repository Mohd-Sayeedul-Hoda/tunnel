-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users(
id SERIAL PRIMARY KEY,
first_name TEXT,
last_name TEXT,
email VARCHAR(50) NOT NULL UNIQUE,
email_verified BOOL DEFAULT FALSE,
created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
