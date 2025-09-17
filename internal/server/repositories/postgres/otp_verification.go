package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/models"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/repositories/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

type otpVerificationRepo struct {
	queries sqlc.Querier
}

func NewOtpVerificationRepo(pool *pgxpool.Pool) (*otpVerificationRepo, error) {
	if pool == nil {
		return nil, errors.New("no pgx pool provided")
	}

	return &otpVerificationRepo{
		queries: sqlc.New(pool),
	}, nil
}

func (o *otpVerificationRepo) CreateOrUpdateOtp(otp *models.OtpVerification) error {

	if otp == nil {
		panic("otp model cannot be nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	err := o.queries.CreateOrUpdateOtp(ctx, sqlc.CreateOrUpdateOtpParams{})
}
