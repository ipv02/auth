package user

import (
	"context"

	"github.com/ipv02/auth/internal/model"
)

// UpdateUser запрос сервесного слоя на обновление данных о пользователе
func (s *service) UpdateUser(ctx context.Context, user *model.UserUpdate) error {
	err := s.userRepository.UpdateUser(ctx, user)

	return err
}
