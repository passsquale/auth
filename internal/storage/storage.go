package storage

import (
	"github.com/go-redis/redis"
	"time"
)

type Redis interface {
	Ping() error
	Close() error
	Get(string) *redis.StringCmd
	Set(string, interface{}, time.Duration) *redis.StatusCmd
}
