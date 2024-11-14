package service

import (
	"context"

	"github.com/ipv02/auth/internal/model"
)

// UserService интерфейс описывающий сервисный слой
type UserService interface {
	CreateUser(ctx context.Context, user *model.UserCreate) (int64, error)
	GetUser(ctx context.Context, id int64) (*model.UserGet, error)
	UpdateUser(ctx context.Context, user *model.UserUpdate) error
	DeleteUser(ctx context.Context, id int64) error
}

// ConsumerService интерфейс описывающий consumer
type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}
