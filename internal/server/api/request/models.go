package request

import (
	"context"
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
	if u.ExpiresAt.Before(time.Now()) {
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

	return v
}

func (u *ForgotPasswordVerify) Valid(ctx context.Context, v *Valid) *Valid {

	v.Check(strings.TrimSpace(u.EmailOtp) != "", "email-otp", "otp should not be empty")
	ValidPassword(v, u.Password)
	return v
}
