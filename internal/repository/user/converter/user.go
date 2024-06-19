package converter

import (
	"github.com/passsquale/auth/internal/model"
	modelRepo "github.com/passsquale/auth/internal/repository/user/model"
)

func ToUserFromRepo(user modelRepo.User) *model.User {
	return &model.User{
		ID:        user.ID,
		Info:      ToUserInfoFromRepo(user.Info),
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func ToUserInfoFromRepo(user modelRepo.UserInfo) model.UserInfo {
	return model.UserInfo{
		Username: user.Name,
		Email:    user.Email,
		Role:     model.UserRole(user.Role),
	}
}
