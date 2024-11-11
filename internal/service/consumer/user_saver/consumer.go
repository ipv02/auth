package user_saver

import (
	"context"

	"github.com/ipv02/auth/internal/client/kafka"
	"github.com/ipv02/auth/internal/repository"
	def "github.com/ipv02/auth/internal/service"
)

var _ def.ConsumerService = (*service)(nil)

type service struct {
	userRepository repository.UserRepository
	consumer       kafka.Consumer
}

// NewService создает и возвращает новый экземпляр сервиса
func NewService(
	userRepository repository.UserRepository,
	consumer kafka.Consumer,
) *service {
	return &service{
		userRepository: userRepository,
		consumer:       consumer,
	}
}

// RunConsumer запускает процесс потребления сообщений в фоновом режиме
func (s *service) RunConsumer(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-s.run(ctx):
			if err != nil {
				return err
			}
		}
	}
}

func (s *service) run(ctx context.Context) <-chan error {
	errChan := make(chan error)

	go func() {
		defer close(errChan)

		errChan <- s.consumer.Consume(ctx, "test-topic", s.UserSaveHandler)
	}()

	return errChan
}
