package user

import (
	"context"
	"log"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/ipv02/auth/pkg/user_v1"
)

// DeleteUser запрос на удаление пользователя.
func (i *Implementation) DeleteUser(ctx context.Context, req *user_v1.DeleteUserRequest) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	err := i.userService.DeleteUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	log.Printf("deleted user: %v", req)

	return &emptypb.Empty{}, nil
}
