package repositories

import "github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/models"

type UserRepo interface {
	// GetById()
	// GetByEmail()

	Insert(*models.User) error
	// Update()
	// Delete()
}
