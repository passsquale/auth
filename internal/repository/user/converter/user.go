package converter

import (
	"github.com/passsquale/auth/internal/repository/user/model"
	desc "github.com/passsquale/auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToUserFromRepo(user model.User) *desc.User {
	var updatedAt *timestamppb.Timestamp
	if user.UpdatedAt.Valid {
		updatedAt = timestamppb.New(user.UpdatedAt.Time)
	}
	return &desc.User{
		Id:        user.ID,
		Info:      ToUserInfoFromRepo(user.UserInfo),
		Password:  user.Password,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: updatedAt,
	}
}

func ToUserInfoFromRepo(user model.UserInfo) *desc.UserInfo {
	return &desc.UserInfo{
		Name:  user.Name,
		Email: user.Email,
		Role:  desc.Role(user.Role),
	}
}
