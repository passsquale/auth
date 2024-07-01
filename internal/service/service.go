package service

import (
	"context"
	"github.com/passsquale/auth/internal/model"
)

type UserService interface {
	Create(ctx context.Context, userCreate *model.UserCreate) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	Update(ctx context.Context, userUpdate *model.UserUpdate) error
	Delete(ctx context.Context, id int64) error
}

type AuthService interface {
	Login(context.Context, model.LoginDTO) (string, error)
	GetRefreshToken(context.Context, string) (string, error)
	GetAccessToken(context.Context, string) (string, error)
}

type AccessService interface {
	Check(ctx context.Context, endpoint string) error
}
