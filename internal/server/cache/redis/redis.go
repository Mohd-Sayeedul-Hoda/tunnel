package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/cache"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/config"

	"github.com/redis/go-redis/v9"
)

var (
	ErrNotFound = errors.New("record not found")
)

type redisRepo struct {
	client *redis.Client
}

func NewRedisCacheRepo(cfg *config.Config) (cache.CacheRepo, error) {
	opts, err := redis.ParseURL(cfg.Cache.DSN)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opts)
	err = client.Ping(context.Background()).Err()
	if err != nil {
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}

	return &redisRepo{
		client: client,
	}, nil
}

func (r *redisRepo) Get(key string) (string, error) {

	value, err := r.client.Get(context.Background(), key).Result()
	if err != nil {
		switch {
		case errors.Is(err, redis.Nil):
			return "", ErrNotFound
		default:
			return "", err
		}
	}
	return value, nil

}

func (r *redisRepo) Set(key, value string, exp time.Duration) error {
	return r.client.Set(context.Background(), key, value, exp).Err()
}

func (r *redisRepo) Delete(key string) (bool, error) {
	cmd := r.client.Del(context.Background(), key)
	if cmd.Err() != nil {
		return false, cmd.Err()
	}
	return cmd.Val() > 0, nil
}
