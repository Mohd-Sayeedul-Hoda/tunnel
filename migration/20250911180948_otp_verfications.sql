-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS otp_verification(
  id SERIAL PRIMARY KEY,
  email VARCHAR(300) NOT NULL,
  otp TEXT NOT NULL,
  type TEXT NOT NULL CHECK (type IN('sign-in', 'email-verification', 'forget-password')),
  attempts INT NOT NULL DEFAULT 0,
  resend_count INT NOT NULL DEFAULT 0,
  used BOOLEAN DEFAULT false,
  is_invalidated BOOLEAN DEFAULT false,
  expires_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(email, type)
);

CREATE INDEX IF NOT EXISTS idx_otp_verifications_email
  ON otp_verification (email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_otp_verifications_email;

DROP TABLE IF  EXISTS  otp_verification;
-- +goose StatementEnd
