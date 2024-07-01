package access

import (
	"github.com/passsquale/auth/internal/service"
	accesspb "github.com/passsquale/auth/pkg/access_v1"
)

type Implementation struct {
	accesspb.UnimplementedAccessV1Server
	service service.AccessService
}

func NewImplementation(serv service.AccessService) *Implementation {
	return &Implementation{
		service: serv,
	}
}
