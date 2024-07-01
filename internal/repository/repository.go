package repository

import (
	"context"
	"github.com/passsquale/auth/internal/model"
	"github.com/passsquale/auth/internal/utils/filter"
)

type UserRepository interface {
	Create(ctx context.Context, userCreate *model.UserCreate) (int64, error)
	Get(ctx context.Context, filter filter.Filter) (*model.User, error)
	Update(ctx context.Context, wrap *model.UserUpdate) error
	Delete(ctx context.Context, id int64) error
}
