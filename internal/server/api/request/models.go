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

type SendOtp struct {
	Email   string `json:"email"`
	OtpType string `json:"type"`
}

type VerfyOtp struct {
	Email    string `json:"email"`
	OtpType  string `json:"type"`
	EmailOtp string `json:"email-otp"`
}

func (u *SendOtp) Valid(ctx context.Context, v *Valid) *Valid {
	ValidEmail(v, u.Email)
	ValidOtpType(v, u.OtpType)
	return v
}

func (u *VerfyOtp) Valid(ctx context.Context, v *Valid) *Valid {

	ValidEmail(v, u.Email)
	ValidOtpType(v, u.OtpType)
	v.Check(strings.TrimSpace(u.EmailOtp) != "", "email-otp", "otp should not be empty")

	return v

}
