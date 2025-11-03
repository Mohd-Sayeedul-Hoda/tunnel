package db

import (
	"context"
	"time"

	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/shared/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func OpenPostgresConn(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.DB.DSN)
	if err != nil {
		return nil, err
	}

	poolConfig.MaxConns = int32(cfg.DB.MaxOpenConn)
	poolConfig.MinConns = int32(cfg.DB.MaxIdealConn)
	poolConfig.MaxConnIdleTime, err = time.ParseDuration(cfg.DB.MaxIdleTime)
	if err != nil {
		return nil, err
	}

	conn, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}

	if err = conn.Ping(ctx); err != nil {
		return nil, err
	}

	return conn, nil
}
