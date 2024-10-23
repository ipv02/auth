package user

import (
	"context"

	"github.com/ipv02/auth/internal/model"
)

func (s *serv) GetUser(ctx context.Context, id int64) (*model.UserGet, error) {
	user, err := s.userRepository.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
