package repositories

import (
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/models"
)

type UserRepo interface {
	GetById(userId int) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Create(user *models.User) error
	Delete(userId int) error
}
