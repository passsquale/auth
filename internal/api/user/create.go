package user

import (
	"context"
	"github.com/passsquale/auth/internal/converter"
	desc "github.com/passsquale/auth/pkg/user_v1"
)

func (i *UserImplementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	id, err := i.userService.Create(ctx, converter.ToUserCreateFromDesc(req))
	if err != nil {
		return nil, err
	}
	return &desc.CreateResponse{Id: id}, nil
}
