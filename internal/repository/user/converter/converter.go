package converter

import (
	"github.com/ipv02/auth/internal/model"
	modelRepo "github.com/ipv02/auth/internal/repository/user/model"
)

// ToUserFromRepo конвертер модели из репо-слоя в модель для сервисного слоя
func ToUserFromRepo(user *modelRepo.User) *model.UserGet {
	return &model.UserGet{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		UserRole:  user.UserRole,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
