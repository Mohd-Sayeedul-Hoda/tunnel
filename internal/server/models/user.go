package models

import "time"

type User struct {
	Id            int       `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Password      string    `json:"-"`
	PasswordHash  []byte    `json:"-"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	EmailVerified bool      `json:"-"`
	IsAdmin       bool      `json:"-"`
}
