-- name: CreateOtp :exec
INSERT INTO otp_verification (email, otp, type, expires_at)
VALUES ($1, $2, $3, $4);

-- name: GetOtp :one
SELECT * FROM otp_verification
WHERE email = $1 AND type = $2
ORDER BY created_at DESC 
LIMIT 1;

-- name: VerifyOtp :exec
UPDATE otp_verification
SET used = true,
    updated_at = NOW()
WHERE id = $1;

-- name: InvalidateOtp :exec
UPDATE otp_verification
SET is_invalidated = true,
    updated_at = NOW()
WHERE id = $1;

-- name: IncreaseOtpAttempt :exec
UPDATE otp_verification
SET attempts = attempts + 1,
    updated_at = NOW()
WHERE id = $1;

-- name: CountOtpsAfterUtcTime :one
SELECT COUNT(*) FROM otp_verification
WHERE created_at > $1 AND email = $2 AND type = $3;

-- name: IncreaseAttemptAndInvalidateOtp :exec
UPDATE otp_verification
SET attempts = attempts + 1,
is_invalidated = true,
updated_at = NOW()
WHERE id = $1;
