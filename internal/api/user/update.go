package user

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/passsquale/auth/internal/converter"
	desc "github.com/passsquale/auth/pkg/user_v1"
)

func (i *UserImplementation) Update(ctx context.Context, req *desc.UpdateRequest) (*empty.Empty, error) {
	err := i.userService.Update(ctx, converter.ToUserUpdateFromDesc(*req))
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}
