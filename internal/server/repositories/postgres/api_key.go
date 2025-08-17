package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/models"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/repositories"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/repositories/sqlc"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type apiKeyRepo struct {
	queries sqlc.Querier
}

func NewAPIKeyRepo(pool *pgxpool.Pool) (repositories.APIRepo, error) {
	if pool == nil {
		return nil, errors.New("no pgx pool provided")
	}

	return &apiKeyRepo{
		queries: sqlc.New(pool),
	}, nil
}

func (a *apiKeyRepo) CreateAPIKey(apiKey *models.APIKey) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var expiredAt pgtype.Timestamptz
	if !apiKey.ExpireAt.IsZero() {
		expiredAt = pgtype.Timestamptz{
			Time:  apiKey.ExpireAt,
			Valid: true,
		}
	} else {
		expiredAt = pgtype.Timestamptz{
			Valid: false,
		}
	}

	createdAPIKey, err := a.queries.CreateAPIKey(ctx, sqlc.CreateAPIKeyParams{
		Name:        apiKey.Name,
		Prefix:      apiKey.Prefix,
		ApiKey:      apiKey.APIKeyHash,
		UserID:      int32(apiKey.UserId),
		Permissions: apiKey.Permissions,
		ExpiresAt:   expiredAt,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return fmt.Errorf("%w: %w", ErrUniqueViolation, err)
			}
		}
		return fmt.Errorf("failed to create API key: %w", err)
	}

	apiKey.Id = int(createdAPIKey.ID)
	apiKey.CreatedAt = createdAPIKey.CreatedAt.Time

	return nil
}

func (a *apiKeyRepo) ListAPIKeys(userId, limit, offset int) ([]models.APIKey, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	keys, err := a.queries.ListAPIKeys(ctx, sqlc.ListAPIKeysParams{
		UserID: int32(userId),
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list api key: %w", err)
	}

	var modelKeys []models.APIKey

	for _, v := range keys {
		modelKey := models.APIKey{
			Id:          int(v.ID),
			Name:        v.Name,
			Prefix:      v.Prefix,
			APIKeyHash:  v.ApiKey,
			UserId:      int(v.UserID),
			ExpireAt:    v.ExpiresAt.Time,
			CreatedAt:   v.CreatedAt.Time,
			Permissions: v.Permissions,
		}

		modelKeys = append(modelKeys, modelKey)
	}

	return modelKeys, nil
}

func (a *apiKeyRepo) DeleteAPIKey(userId, keyId int) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	rows, err := a.queries.DeleteAPIKey(ctx, sqlc.DeleteAPIKeyParams{
		ID:     int32(keyId),
		UserID: int32(userId),
	})
	if err != nil {
		return fmt.Errorf("failed to delete api key: %w", err)
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}
