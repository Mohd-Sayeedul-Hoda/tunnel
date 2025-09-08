package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/models"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/repositories/sqlc"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrUniqueViolation = errors.New("unique key violation in db")
	ErrNotFound        = errors.New("record not found")
)

type userRepo struct {
	queries sqlc.Querier
}

func NewUserRepo(pool *pgxpool.Pool) (*userRepo, error) {
	if pool == nil {
		return nil, errors.New("no pgx pool provided")
	}

	return &userRepo{
		queries: sqlc.New(pool),
	}, nil
}

func (u *userRepo) Create(user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	createdRow, err := u.queries.CreateUser(ctx, sqlc.CreateUserParams{
		Email:        user.Email,
		Name:         user.Name,
		PasswordHash: user.PasswordHash,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return fmt.Errorf("%w: %w", ErrUniqueViolation, err)
			}
		}
		return fmt.Errorf("failed to create user: %w", err)
	}

	user.Id = int(createdRow.ID)
	user.CreatedAt = createdRow.CreatedAt.Time

	return nil
}

func (u *userRepo) Delete(userId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	rows, err := u.queries.DeleteUser(ctx, int32(userId))
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (u *userRepo) GetByEmail(email string) (*models.User, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	dbUser, err := u.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &models.User{
		Id:            int(dbUser.ID),
		Name:          dbUser.Name,
		Email:         dbUser.Email,
		PasswordHash:  dbUser.PasswordHash,
		EmailVerified: dbUser.EmailVerified,
		CreatedAt:     dbUser.CreatedAt.Time,
		UpdatedAt:     dbUser.UpdatedAt.Time,
	}, nil

}

func (u *userRepo) GetById(userId int) (*models.User, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	dbUser, err := u.queries.GetUserById(ctx, int32(userId))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return &models.User{
		Id:            int(dbUser.ID),
		Name:          dbUser.Name,
		Email:         dbUser.Email,
		PasswordHash:  dbUser.PasswordHash,
		EmailVerified: dbUser.EmailVerified,
		CreatedAt:     dbUser.CreatedAt.Time,
		UpdatedAt:     dbUser.UpdatedAt.Time,
	}, nil
}

func (u *userRepo) ListUsers(limit, offset int) ([]models.User, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	dbUsers, err := u.queries.ListUsers(ctx, sqlc.ListUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}

	users := []models.User{}
	for _, dbUser := range dbUsers {

		user := models.User{
			Id:            int(dbUser.ID),
			Name:          dbUser.Name,
			Email:         dbUser.Email,
			PasswordHash:  dbUser.PasswordHash,
			CreatedAt:     dbUser.CreatedAt.Time,
			UpdatedAt:     dbUser.UpdatedAt.Time,
			EmailVerified: dbUser.EmailVerified,
		}
		users = append(users, user)
	}

	return users, nil

}
