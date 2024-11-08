package converter

import (
	"log"

	"github.com/ipv02/auth/internal/model"
	modelRepo "github.com/ipv02/auth/internal/repository/user/pg/model"
)

// ToUserFromRepo конвертер модели из репо-слоя в модель для сервисного слоя
func ToUserFromRepo(user *modelRepo.User) *model.UserGet {
	if user == nil {
		log.Println("ToUserFromRepo: nil user, returning nil response")
		return nil
	}

	return &model.UserGet{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		UserRole:  user.UserRole,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
