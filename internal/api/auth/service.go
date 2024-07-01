package auth

import (
	"github.com/passsquale/auth/internal/service"
	authPb "github.com/passsquale/auth/pkg/auth_v1"
)

type Implementation struct {
	authPb.UnimplementedAuthV1Server
	authService service.AuthService
}

func NewImplementation(service service.AuthService) *Implementation {
	return &Implementation{
		authService: service,
	}
}
