package repository

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/ipv02/auth/pkg/user_v1"
)

type UserRepository interface {
	CreateUser(ctx context.Context, req *user_v1.CreateUserRequest) (*user_v1.CreateUserResponse, error)
	GetUser(ctx context.Context, req *user_v1.GetUserRequest) (*user_v1.GetUserResponse, error)
	UpdateUser(ctx context.Context, req *user_v1.UpdateUserRequest) (*emptypb.Empty, error)
	DeleteUser(ctx context.Context, req *user_v1.DeleteUserRequest) (*emptypb.Empty, error)
}
