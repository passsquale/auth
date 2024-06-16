package user

import (
	"context"
	"github.com/passsquale/auth/internal/converter"
	desc "github.com/passsquale/auth/pkg/user_v1"
)

func (i *UserImplementation) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	user, err := i.userService.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return &desc.GetResponse{
		User: converter.ToUserFromService(*user),
	}, nil
}
