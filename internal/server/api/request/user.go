package request

import "context"

type User struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (u *User) Valid(ctx context.Context) *Valid {
	v := NewValidator()

	validEmail(v, u.Email)
	validPassword(v, u.Password)
	validName(v, u.Name)

	return v
}
