package user

import (
	"github.com/ipv02/auth/internal/client/db"
	"github.com/ipv02/auth/internal/repository"
	userService "github.com/ipv02/auth/internal/service"
)

type service struct {
	userRepository repository.UserRepository
	txManager      db.TxManager
}

// NewService конструктор для создания связи между сервисным слоем и репо слоем
func NewService(userRepository repository.UserRepository, txManger db.TxManager) userService.UserService {
	return &service{
		userRepository: userRepository,
		txManager:      txManger,
	}
}

// NewMockService мок конструктор для создания связи между сервисным слоем и репо слоем
func NewMockService(deps ...interface{}) userService.UserService {
	service := service{}

	for _, v := range deps {
		switch s := v.(type) {
		case repository.UserRepository:
			service.userRepository = s
		}
	}

	return &service
}
