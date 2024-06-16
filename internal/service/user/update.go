package user

import (
	"context"
	"github.com/passsquale/auth/internal/model"
)

func (s *serv) Update(ctx context.Context, userUpdate *model.UserUpdate) error {
	return s.userRepository.Update(ctx, userUpdate)
}
