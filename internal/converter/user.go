package converter

import (
	"github.com/passsquale/auth/internal/model"
	desc "github.com/passsquale/auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToUserFromService(user model.User) *desc.User {
	var updatedAt *timestamppb.Timestamp
	if user.UpdatedAt.Valid {
		updatedAt = timestamppb.New(user.UpdatedAt.Time)
	}
	return &desc.User{
		Id:        user.ID,
		Info:      ToUserInfoFromService(user.Info),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: updatedAt,
	}
}

func ToUserInfoFromService(user model.UserInfo) *desc.UserInfo {
	return &desc.UserInfo{
		Username: user.Username,
		Email:    user.Email,
		Role:     desc.UserRole(user.Role),
	}
}

func ToUserUpdateFromDesc(updateRequest *desc.UpdateRequest) *model.UserUpdate {
	return &model.UserUpdate{
		ID: updateRequest.Id,
		Info: model.UserInfo{
			Username: updateRequest.Info.Username.Value,
			Email:    updateRequest.Info.Email.Value,
			Role:     model.UserRole(updateRequest.Info.Role),
		},
	}
}

func ToUserCreateFromDesc(createRequest *desc.CreateRequest) *model.UserCreate {
	return &model.UserCreate{
		Info: model.UserInfo{
			Username: createRequest.Info.Username,
			Email:    createRequest.Info.Email,
			Role:     model.UserRole(createRequest.Info.Role),
		},
		Password: createRequest.Password,
	}
}
