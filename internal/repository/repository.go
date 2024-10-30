package repository

import (
	"context"

	"github.com/ipv02/auth/internal/model"
)

// UserRepository интерфейс описывающий репо слой
type UserRepository interface {
	CreateUser(ctx context.Context, user *model.UserCreate) (int64, error)
	GetUser(ctx context.Context, id int64) (*model.UserGet, error)
	UpdateUser(ctx context.Context, user *model.UserUpdate) error
	DeleteUser(ctx context.Context, id int64) error
}
