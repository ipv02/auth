package user

import (
	"context"
	"log"

	"github.com/ipv02/auth/internal/converter"
	"github.com/ipv02/auth/pkg/user_v1"
)

// CreateUser - запрос создает нового пользователя.
func (i *Implementation) CreateUser(ctx context.Context, req *user_v1.CreateUserRequest) (*user_v1.CreateUserResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	id, err := i.userService.CreateUser(ctx, converter.ToUserCreateFromReq(req))
	if err != nil {
		return nil, err
	}

	log.Printf("created user: %v", id)

	return &user_v1.CreateUserResponse{
		Id: id,
	}, nil
}
