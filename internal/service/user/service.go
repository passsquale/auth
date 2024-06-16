package user

import (
	"github.com/passsquale/auth/internal/client/db"
	"github.com/passsquale/auth/internal/repository"
)

type serv struct {
	userRepository repository.UserRepository
	txManager      db.TxManager
}

func NewUserService(userRepo repository.UserRepository, txManager db.TxManager) repository.UserRepository {
	return &serv{
		userRepository: userRepo,
		txManager:      txManager,
	}
}
