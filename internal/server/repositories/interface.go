package repositories

import (
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/models"
)

type UserRepo interface {
	ListUsers(limit, offset int32) ([]models.User, error)
	GetById(userId int) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Create(user *models.User) error
	Delete(userId int) error
}

type APIRepo interface {
	CreateAPIKey(apiKey *models.APIKey) error
	ListAPIKeys(userId, limit, offset int32) ([]models.APIKey, error)
	DeleteAPIKey(keyId, userId int32) error
}
