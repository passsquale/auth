package auth

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/passsquale/auth/internal/model"
	"github.com/passsquale/auth/internal/utils"
	"github.com/passsquale/auth/internal/utils/filter"
)

func (s *serv) Login(ctx context.Context, req model.LoginDTO) (string, error) {
	conditions := filter.MakeFilter(filter.Condition{
		Key:   model.UserNameFieldCode,
		Value: req.Username,
	})

	user, err := s.userRepo.Get(ctx, conditions)
	if err != nil {
		return "", err
	}

	if !utils.VerifyPassword(user.Password, req.Password) {
		return "", errors.New("authentication failed")
	}

	token, err := utils.GenerateToken(user.Info, s.jwtConfig.RefreshSecretKey(), s.jwtConfig.RefreshExpirationTime())
	if err != nil {
		return "", err
	}

	infoJSON, err := json.Marshal(user.Info)
	if err != nil {
		return "", err
	}
	res := s.redis.Set(user.Info.Username, infoJSON, 0)
	if res.Err() != nil {
		return "", err
	}

	return token, nil
}
