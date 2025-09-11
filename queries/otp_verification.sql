-- name: CreateOrUpdateOtp :exec
INSERT INTO otp_verification (email, otp, type, expires_at, resend_count, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (email, type) DO UPDATE
SET otp = EXCLUDED.otp,
    expires_at = EXCLUDED.expires_at,
    resend_count = otp_verification.resend_count + 1,
    attempts = 0,
    used = false,
    is_invalidated = false,
    updated_at = NOW();

-- name: GetOtp :one
SELECT * FROM otp_verification
WHERE email = $1 AND type = $2;

-- name: VerifyOtp :exec
UPDATE otp_verification
SET used = true,
    updated_at = NOW()
WHERE email = $1 AND type = $2;

-- name: InvalidateOtp :exec
UPDATE otp_verification
SET is_invalidated = true,
    updated_at = NOW()
WHERE email = $1 AND type = $2;

-- name: IncreaseOtpAttempt :exec
UPDATE otp_verification
SET attempts = attempts + 1,
    updated_at = NOW()
WHERE email = $1 AND type = $2;
