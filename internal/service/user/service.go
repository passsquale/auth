package user

import (
	"github.com/passsquale/auth/internal/repository"
)

type serv struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) repository.UserRepository {
	return &serv{userRepository: userRepo}
}
