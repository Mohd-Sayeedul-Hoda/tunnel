package repositories

import (
	"time"

	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/models"
)

type UserRepo interface {
	ListUsers(limit, offset int) ([]models.User, error)
	GetById(userId int) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Create(user *models.User) error
	Delete(userId int) error
}

type APIRepo interface {
	CreateAPIKey(apiKey *models.APIKey) error
	ListAPIKeys(userId, limit, offset int) ([]models.APIKey, error)
	DeleteAPIKey(userId, keyId int) error
}

type OtpVerificationRepo interface {
	CreateOtp(email, otp string, typeOfOtp models.OtpType, expiersAt time.Time) error
	GetOtp(email string, otpType models.OtpType) (*models.OtpVerification, error)
	VerifyOtp(id int) error
	InvalidateOtp(id int) error
	IncreaseOtpAttempt(id int) error
	CountOtpsAfterUtcTime(email string, otpType models.OtpType, after time.Time) (int, error)
	IncreaseAttemptAndInvalidateOtp(id int) error
}
