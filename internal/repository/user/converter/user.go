package converter

import (
	"github.com/passsquale/auth/internal/repository/user/model"
	desc "github.com/passsquale/auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToUserFromRepo(user *model.User) *desc.GetResponse {
	var updatedAt *timestamppb.Timestamp
	if user.UpdatedAt.Valid {
		updatedAt = timestamppb.New(user.UpdatedAt.Time)
	}
	return &desc.GetResponse{
		Id:        user.ID,
		Info:      ToUserInfoFromRepo(user.Info),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: updatedAt,
	}
}
func ToUserInfoFromRepo(info *model.Info) *desc.UserInfo {
	return &desc.UserInfo{
		Name:  info.Name,
		Email: info.Email,
		Role:  desc.Role(info.Role),
	}
}
