package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/models"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/repositories/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type otpVerificationRepo struct {
	queries sqlc.Querier
}

func NewEmailOtpRepo(pool *pgxpool.Pool) (*otpVerificationRepo, error) {
	if pool == nil {
		return nil, errors.New("no pgx pool provided")
	}

	return &otpVerificationRepo{
		queries: sqlc.New(pool),
	}, nil
}

func (o *otpVerificationRepo) CreateOtp(email, otp string, typeOfOtp models.OtpType, expiersAt time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	err := o.queries.CreateOtp(ctx, sqlc.CreateOtpParams{
		Email: email,
		Otp:   otp,
		Type:  string(typeOfOtp),
		ExpiresAt: pgtype.Timestamptz{
			Time:  expiersAt,
			Valid: true,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create email otp: %w", err)
	}

	return nil
}

func (o *otpVerificationRepo) GetOtp(email string, otpType models.OtpType) (*models.OtpVerification, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	otpModel, err := o.queries.GetOtp(ctx, sqlc.GetOtpParams{
		Email: email,
		Type:  string(otpType),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get email otp by id: %w", err)
	}

	return &models.OtpVerification{
		Id:            int(otpModel.ID),
		Email:         otpModel.Email,
		EmailOtp:      otpModel.Otp,
		Type:          models.OtpType(otpModel.Type),
		ExpiresAt:     otpModel.ExpiresAt.Time,
		Attempts:      int(otpModel.Attempts),
		Used:          otpModel.Used.Bool,
		IsInvalidated: otpModel.IsInvalidated.Bool,
		CreatedAt:     otpModel.CreatedAt.Time,
		UpdatedAt:     otpModel.UpdatedAt.Time,
	}, nil
}

func (o *otpVerificationRepo) VerifyOtp(id int) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	err := o.queries.VerifyOtp(ctx, int32(id))
	if err != nil {
		return fmt.Errorf("failed to verify email otp: %w", err)
	}
	return nil
}

func (o *otpVerificationRepo) InvalidateOtp(id int) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	err := o.queries.InvalidateOtp(ctx, int32(id))
	if err != nil {
		return fmt.Errorf("failed to invalidate otp: %w", err)
	}
	return nil
}

func (o *otpVerificationRepo) IncreaseOtpAttempt(id int) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	err := o.queries.IncreaseOtpAttempt(ctx, int32(id))
	if err != nil {
		return fmt.Errorf("failed to increase otp attempt: %w", err)
	}
	return nil
}

func (o *otpVerificationRepo) CountOtpsAfterUtcTime(email string, otpType models.OtpType, after time.Time) (int, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	count, err := o.queries.CountOtpsAfterUtcTime(ctx, sqlc.CountOtpsAfterUtcTimeParams{
		CreatedAt: pgtype.Timestamptz{Time: after, Valid: true},
		Email:     email,
		Type:      string(otpType),
	})
	if err != nil {
		return 0, fmt.Errorf("failed to count number of email otp after time: %w", err)
	}

	return int(count), nil
}

func (o *otpVerificationRepo) IncreaseAttemptAndInvalidateOtp(id int) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	err := o.queries.IncreaseAttemptAndInvalidateOtp(ctx, int32(id))
	if err != nil {
		return fmt.Errorf("failed to increase attempt and invalided otp: %w", err)
	}

	return nil
}
