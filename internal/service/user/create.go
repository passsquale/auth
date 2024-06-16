package user

import (
	"context"
	"github.com/passsquale/auth/internal/model"
)

func (s *serv) Create(ctx context.Context, userCreate *model.UserCreate) (int64, error) {
	var id int64
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		id, errTx = s.userRepository.Create(ctx, userCreate)
		if errTx != nil {
			return errTx
		}
		return nil
	})

	if err != nil {
		return 0, err
	}
	return id, nil
}
