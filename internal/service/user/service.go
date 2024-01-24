package user

import (
	"github.com/passsquale/auth/internal/repository"
	"github.com/passsquale/auth/internal/service"
)

type serv struct {
	userRepository repository.UserRepository
}

func NewService(userRepository repository.UserRepository) service.UserService {
	return &serv{
		userRepository: userRepository,
	}
}
