package converter

import (
	"database/sql"
	"log"
	"time"

	"github.com/ipv02/auth/internal/model"
	modelRepo "github.com/ipv02/auth/internal/repository/user/redis/model"
)

// ToUserFromRepo конвертер модели из репо-слоя в модель для сервисного слоя
func ToUserFromRepo(user *modelRepo.User) *model.UserGet {
	if user == nil {
		log.Println("ToUserFromRepo: nil user, returning nil response")
		return nil
	}

	var updateAt sql.NullTime
	if user.UpdatedAtNs != nil {
		updateAt = sql.NullTime{
			Time:  time.Unix(0, *user.UpdatedAtNs),
			Valid: true,
		}
	}

	return &model.UserGet{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		UserRole:  user.UserRole,
		CreatedAt: time.Unix(0, user.CreatedAtNs),
		UpdatedAt: updateAt,
	}
}
