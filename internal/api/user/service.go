package user

import (
	"github.com/ipv02/auth/internal/service"
	"github.com/ipv02/auth/pkg/user_v1"
)

// Implementation структура описавающая сервер
type Implementation struct {
	user_v1.UnimplementedUserV1Server
	userService service.UserService
}

// NewImplementation конструктор создает реализацию сервера и связывает ее с бизнес-логиклй
func NewImplementation(userService service.UserService) *Implementation {
	return &Implementation{
		userService: userService,
	}
}
