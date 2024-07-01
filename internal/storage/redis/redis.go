package redis

import (
	"github.com/go-redis/redis"
	"github.com/passsquale/auth/internal/storage"
	"time"
)

type redisStorage struct {
	client *redis.Client
}

func NewRedisConnection(opts *redis.Options) (storage.Redis, error) {
	client := redis.NewClient(opts)

	return &redisStorage{
		client: client,
	}, nil
}

func (r *redisStorage) Ping() error {
	return r.client.Ping().Err()
}

func (r *redisStorage) Close() error {
	return r.client.Close()
}

func (r *redisStorage) Get(key string) *redis.StringCmd {
	return r.client.Get(key)
}

func (r *redisStorage) Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return r.client.Set(key, value, expiration)
}
