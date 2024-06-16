package user

import (
	"context"
	"github.com/passsquale/auth/internal/repository"
)

type serv struct {
	userRepository repository.UserRepository
}

func NewUserService(ctx context.Context, userRepo repository.UserRepository) repository.UserRepository {
	return &serv{userRepository: userRepo}
}
