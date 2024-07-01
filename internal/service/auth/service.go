package auth

import (
	"github.com/passsquale/auth/internal/config"
	"github.com/passsquale/auth/internal/repository"
	"github.com/passsquale/auth/internal/service"
	"github.com/passsquale/auth/internal/storage"
)

type serv struct {
	redis     storage.Redis
	userRepo  repository.UserRepository
	jwtConfig config.JWTConfig
}

func NewService(redis storage.Redis, repo repository.UserRepository, jwtConfig config.JWTConfig) service.AuthService {
	return &serv{
		redis:     redis,
		userRepo:  repo,
		jwtConfig: jwtConfig,
	}
}
