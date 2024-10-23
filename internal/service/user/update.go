package user

import (
	"context"

	"github.com/ipv02/auth/internal/model"
)

func (s *serv) UpdateUser(ctx context.Context, user *model.UserUpdate) error {
	err := s.userRepository.UpdateUser(ctx, user)

	return err
}
