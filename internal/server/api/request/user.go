package request

import (
	"context"
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
