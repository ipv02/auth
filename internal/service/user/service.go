package user

import (
	"github.com/ipv02/auth/internal/repository"
	"github.com/ipv02/auth/internal/service"
)

type serv struct {
	userRepository repository.UserRepository
}

// NewService конструктор для создания связи между сервисным слоем и репо слоем
func NewService(userRepository repository.UserRepository) service.UserService {
	return &serv{userRepository: userRepository}
}
