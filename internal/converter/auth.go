package converter

import (
	"github.com/passsquale/auth/internal/model"
	authPb "github.com/passsquale/auth/pkg/auth_v1"
)

func AuthProtoToAuthDTO(req *authPb.LoginRequest) model.LoginDTO {
	return model.LoginDTO{
		Password: req.GetPassword(),
		Username: req.GetUsername(),
	}
}
