package user

import (
	"github.com/passsquale/auth/internal/client/db"
	"github.com/passsquale/auth/internal/repository"
	"github.com/passsquale/auth/internal/service"
	"github.com/passsquale/auth/internal/storage"
)

type serv struct {
	userRepository repository.UserRepository
	txManager      db.TxManager
	storage        storage.Redis
}

func NewUserService(userRepo repository.UserRepository, txManager db.TxManager, storage storage.Redis) service.UserService {
	return &serv{
		userRepository: userRepo,
		txManager:      txManager,
		storage:        storage,
	}
}
