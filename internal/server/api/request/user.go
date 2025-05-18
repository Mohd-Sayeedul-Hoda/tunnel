package request

import "context"

type User struct {
	Email    string
	Name     string
	Password string
}

func (u *User) Valid(ctx context.Context) map[string]string {
	problem := make(map[string]string)

	validEmail(u.Email, problem)
	validPassword(u.Password, problem)
	validName(u.Name, problem)

	return problem
}
