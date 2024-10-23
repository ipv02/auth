package user

import (
	"context"

	"github.com/ipv02/auth/internal/model"
)

func (s *serv) CreateUser(ctx context.Context, user *model.UserCreate) (int64, error) {
	id, err := s.userRepository.CreateUser(ctx, user)
	if err != nil {
		return 0, err
	}

	return id, nil
}
