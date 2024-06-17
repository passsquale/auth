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
		Info:      ToUserInfoFromService(user.UserInfo),
		Password:  user.Password,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: updatedAt,
	}
}

func ToUserInfoFromService(user model.UserInfo) *desc.UserInfo {
	return &desc.UserInfo{
		Name:  user.Name,
		Email: user.Email,
		Role:  desc.Role(user.Role),
	}
}

func ToUserUpdateFromDesc(updateRequest *desc.UpdateRequest) *model.UserUpdate {
	return &model.UserUpdate{
		ID: updateRequest.Wrap.Id,
		UserInfo: model.UserInfo{
			Name:  updateRequest.Wrap.Name.Value,
			Email: updateRequest.Wrap.Email.Value,
			Role:  model.Role(*updateRequest.Wrap.Role.Enum()),
		},
	}
}

func ToUserCreateFromDesc(createRequest *desc.CreateRequest) *model.UserCreate {
	return &model.UserCreate{
		UserInfo: model.UserInfo{
			Name:  createRequest.Info.Name,
			Email: createRequest.Info.Email,
			Role:  model.Role(createRequest.Info.Role),
		},
		Password: createRequest.Password,
	}
}
