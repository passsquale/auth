package converter

import (
	"github.com/passsquale/auth/internal/model"
	modelRepo "github.com/passsquale/auth/internal/repository/user/model"
)

func ToUserFromRepo(user modelRepo.User) *model.User {
	return &model.User{
		ID:        user.ID,
		UserInfo:  *ToUserInfoFromRepo(user.UserInfo),
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func ToUserInfoFromRepo(user modelRepo.UserInfo) *model.UserInfo {
	return &model.UserInfo{
		Name:  user.Name,
		Email: user.Email,
		Role:  model.Role(user.Role),
	}
}
