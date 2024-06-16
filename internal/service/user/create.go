package user

import (
	"context"
	"github.com/passsquale/auth/internal/model"
)

func (s *serv) Create(ctx context.Context, userCreate *model.UserCreate) (int64, error) {
	return s.userRepository.Create(ctx, userCreate)
}
