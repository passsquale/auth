package user

import (
	"github.com/passsquale/auth/internal/service"
	desc "github.com/passsquale/auth/pkg/user_v1"
)

type UserImplementation struct {
	desc.UnimplementedUserV1Server
	userService service.UserService
}

func NewUserImplementation(userService service.UserService) *UserImplementation {
	return &UserImplementation{
		userService: userService,
	}
}
