package user

import (
	"context"
	"github.com/passsquale/auth/internal/model"
	"github.com/passsquale/auth/internal/utils"
	"github.com/passsquale/auth/internal/utils/filter"
	"github.com/pkg/errors"
)

func (s *serv) Create(ctx context.Context, userCreate *model.UserCreate) (int64, error) {
	conditions := filter.MakeFilter(filter.Condition{
		Key:   model.UserNameFieldCode,
		Value: userCreate.Info.Username,
	})

	user, err := s.userRepository.Get(ctx, conditions)
	if err != nil {
		return 0, err
	}

	if user != nil {
		return 0, errors.Errorf(`user with username "%s" already exist`, userCreate.Info.Username)
	}

	hashedPassword, err := utils.HashPassword(userCreate.Password)
	if err != nil {
		return 0, errors.Errorf("failed hash password: %v", err)
	}

	userCreate.Password = hashedPassword

	var id int64

	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error

		id, errTx = s.userRepository.Create(ctx, userCreate)
		if errTx != nil {
			return errTx
		}

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
