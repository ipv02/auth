package user

import (
	"github.com/ipv02/auth/internal/client/db"
	"github.com/ipv02/auth/internal/repository"
	"github.com/ipv02/auth/internal/service"
)

type serv struct {
	userRepository repository.UserRepository
	txManager      db.TxManager
}

// NewService конструктор для создания связи между сервисным слоем и репо слоем
func NewService(userRepository repository.UserRepository, txManger db.TxManager) service.UserService {
	return &serv{
		userRepository: userRepository,
		txManager:      txManger,
	}
}
