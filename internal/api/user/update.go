package user

import (
	"context"
	"log"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/ipv02/auth/internal/converter"
	"github.com/ipv02/auth/pkg/user_v1"
)

// UpdateUser запрос на обновление данных о пользователе.
func (i *Implementation) UpdateUser(ctx context.Context, req *user_v1.UpdateUserRequest) (*emptypb.Empty, error) {
	err := i.userService.UpdateUser(ctx, converter.ToUserUpdateFromReq(req))
	if err != nil {
		return nil, err
	}

	log.Printf("updated user: %v", req)

	return &emptypb.Empty{}, nil
}
