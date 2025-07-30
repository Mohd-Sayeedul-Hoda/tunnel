package request

import "context"

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
