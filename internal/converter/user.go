package converter

import (
	"github.com/passsquale/auth/internal/model"
	desc "github.com/passsquale/auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToUserFromService(user *model.User) *desc.GetResponse {
	var updatedAt *timestamppb.Timestamp
	if user.UpdatedAt.Valid {
		updatedAt = timestamppb.New(user.UpdatedAt.Time)
	}
	return &desc.GetResponse{
		Id:        user.ID,
		Info:      ToUserInfoFromService(user.Info),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: updatedAt,
	}
}
func ToUserInfoFromService(info *model.UserInfo) *desc.UserInfo {
	return &desc.UserInfo{
		Name:  info.Name,
		Email: info.Email,
		Role:  desc.Role(info.Role),
	}
}

func ToUserInfoFromDesc(info *desc.CreateRequest) *model.UserInfo {
	return &model.UserInfo{
		Name:  info.Info.Name,
		Email: info.Info.Email,
		Role:  model.Role(info.Info.Role),
	}
}
