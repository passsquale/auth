package config

import (
	"github.com/joho/godotenv"
	"github.com/passsquale/auth/internal/model"
	"time"
)

func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}
	return nil
}

type GRPCConfig interface {
	Address() string
}

type HTTPConfig interface {
	Address() string
}

type PGConfig interface {
	DSN() string
}

type RedisConfig interface {
	Address() string
	Password() string
	RoutesAccesses() map[string][]model.UserRole
}

type JWTConfig interface {
	RefreshSecretKey() []byte
	RefreshExpirationTime() time.Duration
	AccessSecretKey() []byte
	AccessExpirationTime() time.Duration
}

type SwaggerConfig interface {
	Address() string
}
