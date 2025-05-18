package models

import "time"

type User struct {
	Id           string
	Name         string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}
