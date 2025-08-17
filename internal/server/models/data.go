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

type APIKey struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Prefix      string    `json:"prefix"`
	APIkeyToken string    `json:"api_key_token,omitempty"`
	APIKeyHash  string    `json:"-"`
	UserId      int       `json:"user_id"`
	ExpireAt    time.Time `json:"expire_at"`
	CreatedAt   time.Time `json:"created_at"`
	Permissions []string  `json:"permission,omitempty"`
}
