package access

import (
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes/empty"
	accesspb "github.com/passsquale/auth/pkg/access_v1"
)

func (i *Implementation) Check(ctx context.Context, req *accesspb.CheckRequest) (*empty.Empty, error) {

	err := i.service.Check(ctx, req.GetEndpointAddress())
	if err != nil {
		return nil, errors.New("access denied")
	}

	return &empty.Empty{}, nil
}
