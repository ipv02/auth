package user

import (
	"context"
	"log"

	"github.com/ipv02/auth/internal/converter"
	"github.com/ipv02/auth/pkg/user_v1"
)

// GetUser запроолс получения информации о пользователе.
func (i *Implementation) GetUser(ctx context.Context, req *user_v1.GetUserRequest) (*user_v1.GetUserResponse, error) {
	if err := req.ValidateRequest(); err != nil {
		return nil, err
	}

	userObj, err := i.userService.GetUser(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	log.Printf("get user: %v", userObj)

	return converter.ToUserFromService(userObj), nil
}
