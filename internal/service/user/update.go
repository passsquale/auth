package user

import (
	"context"
	"github.com/passsquale/auth/internal/model"
)

func (s *serv) Update(ctx context.Context, userUpdate *model.UserUpdate) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error

		errTx = s.userRepository.Update(ctx, userUpdate)
		if errTx != nil {
			return errTx
		}

		if errTx != nil {
			return errTx
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
