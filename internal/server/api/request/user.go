package request

import "context"

type User struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (u *User) Valid(ctx context.Context) map[string]string {
	problem := make(map[string]string)

	validEmail(u.Email, problem)
	validPassword(u.Password, problem)
	validName(u.Name, problem)

	return problem
}
