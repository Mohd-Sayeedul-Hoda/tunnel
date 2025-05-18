package postgres

import (
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/repositories"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) repositories.UserRepo {
	return &userRepo{
		pool: pool,
	}
}

func (u *userRepo) Insert() {}
