package repository

import (
	"context"
	desc "github.com/passsquale/auth/pkg/user_v1"
)

type UserRepository interface {
	Create(ctx context.Context, info *desc.CreateRequest) (int64, error)
	Get(ctx context.Context, id int64) (*desc.GetResponse, error)
}
