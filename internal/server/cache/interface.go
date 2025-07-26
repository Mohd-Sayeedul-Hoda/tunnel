package cache

import "time"

type CacheRepo interface {
	Get(key string) (string, error)
	Set(key string, value any, exp time.Duration) error
	Delete(key string) (bool, error)
}
