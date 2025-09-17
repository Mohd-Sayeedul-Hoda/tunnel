-- name: CreateOrUpdateOtp :exec
INSERT INTO otp_verification (email, otp, type, expires_at, updated_at)
VALUES ($1, $2, $3, $4, $5);

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
