package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ipv02/auth/internal/repository/user/model"
	"github.com/ipv02/auth/pkg/user_v1"
)

func ToUserFromRepo(user *model.User) *user_v1.GetUserResponse {
	var updatedAt *timestamppb.Timestamp
	if user.UpdatedAt.Valid {
		updatedAt = timestamppb.New(user.UpdatedAt.Time)
	}

	return &user_v1.GetUserResponse{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user_v1.UserRole(user.UserRole),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: updatedAt,
	}
}
