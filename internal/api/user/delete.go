package user

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	desc "github.com/passsquale/auth/pkg/user_v1"
)

func (i *UserImplementation) Delete(ctx context.Context, req *desc.DeleteRequest) (*empty.Empty, error) {
	err := i.userService.Delete(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, err
}
