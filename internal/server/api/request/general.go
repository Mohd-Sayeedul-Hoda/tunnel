package request

import (
	"regexp"
)

func check(ok bool, key, message string, problems map[string]string) {
	if !ok {
		if _, exists := problems[key]; !exists {
			problems[key] = message
		}
	}
}

func validEmail(email string, problem map[string]string) {
	check(email != "", "email", "email must not be empty", problem)

	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	check(emailRegex.MatchString(email), "email", "email must be valid address", problem)
}

func validPassword(password string, problem map[string]string) {
	check(password != "", "password", "password should not be empty string", problem)
	check(len(password) >= 8, "password", "lenght of password should be greater or equal to 8 character", problem)
	check(len(password) <= 50, "password", "lenght of password should be smaller then 50 character", problem)
}

func validName(name string, problem map[string]string) {
	check(name != "", "name", "name should not be empty string", problem)
	check(len(name) <= 300, "name", "name should be less then 300 character", problem)
}
