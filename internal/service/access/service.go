package access

import (
	"github.com/passsquale/auth/internal/config"
	"github.com/passsquale/auth/internal/service"
	"github.com/passsquale/auth/internal/utils"
)

type serv struct {
	jwtConfig     config.JWTConfig
	accessChecker utils.AccessChecker
}

func NewService(jwtConfig config.JWTConfig, checker utils.AccessChecker) service.AccessService {
	return &serv{
		jwtConfig:     jwtConfig,
		accessChecker: checker,
	}
}
