package request

import (
	"regexp"
)

func ValidEmail(v *Valid, email string) {
	v.Check(email != "", "email", "email must not be empty")

	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	v.Check(emailRegex.MatchString(email), "email", "email must be valid address")
}

func ValidPassword(v *Valid, password string) {
	v.Check(password != "", "password", "password should not be empty string")
	v.Check(len(password) >= 8, "password", "lenght of password should be greater or equal to 8 character")
	v.Check(len(password) <= 50, "password", "lenght of password should be smaller then 50 character")
}

func ValidName(v *Valid, name string) {
	v.Check(name != "", "name", "name should not be empty string")
	v.Check(len(name) <= 300, "name", "name should be less then 300 character")
}

func ValidOtpType(v *Valid, otpType string) {
	validTypes := map[string]bool{
		"email-verification": true,
		"forget-password":    true,
	}
	v.Check(validTypes[otpType], "type", "invalid OTP type. Valid types are: email-verification, forget-password")
}
