package auth

import (
	"context"
	"github.com/passsquale/auth/internal/converter"
	"github.com/passsquale/auth/pkg/auth_v1"
	"github.com/pkg/errors"
)

func (i *Implementation) Login(ctx context.Context, req *auth_v1.LoginRequest) (*auth_v1.LoginResponse, error) {
	refreshToken, err := i.authService.Login(ctx, converter.AuthProtoToAuthDTO(req))
	if err != nil {
		return nil, errors.Errorf("authentification error: %s", err)
	}

	return &auth_v1.LoginResponse{
		RefreshToken: refreshToken,
	}, err
}
