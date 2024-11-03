package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ipv02/auth/internal/model"
	"github.com/ipv02/auth/pkg/user_v1"
)

// ToUserFromService конвертер модели бизнес-логики в протомодель
func ToUserFromService(user *model.UserGet) *user_v1.GetUserResponse {
	if user == nil {
		return nil
	}

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

// ToUserCreateFromReq конвертер протомодели в модель бизнес-логики
func ToUserCreateFromReq(user *user_v1.CreateUserRequest) *model.UserCreate {
	if user == nil {
		return nil
	}

	return &model.UserCreate{
		Name:            user.Name,
		Email:           user.Email,
		Password:        user.Password,
		PasswordConfirm: user.PasswordConfirm,
		Role:            int32(user.Role),
	}
}

// ToUserUpdateFromReq конвертер протомодели в модель бизнес-логики
func ToUserUpdateFromReq(user *user_v1.UpdateUserRequest) *model.UserUpdate {
	if user == nil {
		return nil
	}

	name := user.Name.GetValue()
	email := user.Email.GetValue()

	return &model.UserUpdate{
		ID:    user.Id,
		Name:  &name,
		Email: &email,
		Role:  int32(user.Role),
	}
}
