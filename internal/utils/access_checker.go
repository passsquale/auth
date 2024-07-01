package utils

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis"
	"github.com/passsquale/auth/internal/config"
	"github.com/passsquale/auth/internal/model"
	"github.com/passsquale/auth/internal/repository"
	"github.com/passsquale/auth/internal/storage"
	"github.com/passsquale/auth/internal/utils/filter"
)

type AccessChecker interface {
	AccessCheck(ctx context.Context, token string, endpoint string) (bool, error)
}

type routeAccessChecker struct {
	jwtConfig config.JWTConfig
	redis     storage.Redis
	userRepo  repository.UserRepository
}

func NewRouteAccessChecker(jwtCfg config.JWTConfig, redis storage.Redis, repo repository.UserRepository) AccessChecker {
	return &routeAccessChecker{
		jwtConfig: jwtCfg,
		redis:     redis,
		userRepo:  repo,
	}
}

func (r *routeAccessChecker) AccessCheck(ctx context.Context, token string, endpoint string) (bool, error) {
	claims, err := VerifyToken(token, r.jwtConfig.AccessSecretKey())
	if err != nil {
		return false, err
	}

	var info *model.UserInfo

	res, err := r.redis.Get(claims.Username).Result()
	if errors.Is(err, redis.Nil) {
		conditions := filter.MakeFilter(filter.Condition{
			Key:   model.UserNameFieldCode,
			Value: claims.Username,
		})

		user, errRep := r.userRepo.Get(ctx, conditions)
		if errRep != nil {
			return false, errRep
		}

		info = &user.Info
	}
	if err != nil {
		return false, err
	}

	if info == nil {
		err = json.Unmarshal([]byte(res), &info)
		if err != nil {
			return false, err
		}
	}

	if info.Role == model.ADMIN && claims.Role == model.ADMIN {
		return true, nil
	}

	res, err = r.redis.Get(endpoint).Result()
	if errors.Is(err, redis.Nil) {
		return true, nil
	}
	if err != nil {
		return false, err
	}

	var roles []model.UserRole
	err = json.Unmarshal([]byte(res), &roles)
	if err != nil {
		return false, err
	}

	for _, role := range roles {
		if role == info.Role {
			return true, nil
		}
	}

	return false, nil
}
