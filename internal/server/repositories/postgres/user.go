package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/models"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/repositories"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email address")
	ErrInsertingRow   = errors.New("error while inserting row in db")
)

type userRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) (repositories.UserRepo, error) {
	if pool == nil {
		return nil, errors.New("no pgx pool provided")
	}

	return &userRepo{
		pool: pool,
	}, nil
}

func (u *userRepo) Insert(user *models.User) error {
	query := `INSERT INTO users (email, name, password_hash) VALUES ($1, $2, $3) RETURNING id, created_at`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	err := u.pool.QueryRow(ctx, query,
		user.Email,
		user.Name,
		user.PasswordHash,
	).Scan(user.Id, user.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return fmt.Errorf("%w %w", ErrDuplicateEmail, err)
			}
			return fmt.Errorf("%w %w", ErrInsertingRow, err)
		}
	}

	return nil
}
