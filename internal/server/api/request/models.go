package request

import (
	"context"
	"encoding/base64"
	"strings"
	"time"
)

type User struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (u *User) Valid(ctx context.Context, v *Valid) *Valid {

	ValidEmail(v, u.Email)
	ValidPassword(v, u.Password)
	ValidName(v, u.Name)

	return v
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *Login) Valid(ctx context.Context, v *Valid) *Valid {

	ValidEmail(v, u.Email)
	ValidPassword(v, u.Password)

	return v
}

type APIKeys struct {
	Name      string    `json:"name"`
	ExpiresAt time.Time `json:"expires_at,omitzero"`
}

func (u *APIKeys) Valid(ctx context.Context, v *Valid) *Valid {
	ValidName(v, u.Name)
	if !u.ExpiresAt.IsZero() && u.ExpiresAt.Before(time.Now()) {
		v.AddError("expires_at", "must be in the future")
	}

	return v
}

type BaseEmail struct {
	Email string `json:"email"`
}

type VerifyUserOTP struct {
	Email    string `json:"email"`
	EmailOtp string `json:"otp"`
}

type ForgotPasswordVerify struct {
	EmailOtp string `json:"otp"`
	Password string `json:"password"`
}

func (u *BaseEmail) Valid(ctx context.Context, v *Valid) *Valid {
	ValidEmail(v, u.Email)
	return v
}

func (u *VerifyUserOTP) Valid(ctx context.Context, v *Valid) *Valid {

	ValidEmail(v, u.Email)
	v.Check(strings.TrimSpace(u.EmailOtp) != "", "email-otp", "otp should not be empty")
	v.Check(len(u.EmailOtp) <= 200, "email-otp", "otp too long")
	return v
}

func (u *ForgotPasswordVerify) Valid(ctx context.Context, v *Valid) *Valid {

	v.Check(strings.TrimSpace(u.EmailOtp) != "", "email-otp", "otp should not be empty")
	v.Check(len(u.EmailOtp) <= 200, "email-otp", "otp too long")

	if v.Valid() {
		decodedBytes, err := base64.StdEncoding.DecodeString(u.EmailOtp)
		if err != nil {
			v.AddError("email-otp", "otp format invalid")
		} else {
			decoded := string(decodedBytes)

			if !strings.Contains(decoded, "|") {
				v.AddError("email-otp", "otp format invalid")
			} else {
				parts := strings.SplitN(decoded, "|", 2)

				if len(parts) != 2 {
					v.AddError("email-otp", "otp format invalid")
				}
				email := strings.TrimSpace(parts[0])
				token := parts[1]

				ValidEmail(v, email)

				if len(token) != 32 {
					v.AddError("email-otp", "otp format invalid")
				}
				ValidAlphanumeric(v, token, "email-otp")
			}
		}
	}

	ValidPassword(v, u.Password)
	return v
}

type VerifyAPIKey struct {
	Key string `json:"api_key"`
}

func (u *VerifyAPIKey) Valid(ctx context.Context, v *Valid) *Valid {
	v.Check(u.Key != "", "api_key", "api key should not be empty")
	v.Check(len(u.Key) <= 200, "api_key", "api key too long")

	return v
}
