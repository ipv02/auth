package kafka

import (
	"context"

	"github.com/ipv02/auth/internal/client/kafka/consumer"
)

// Consumer определяет интерфейс для потребителя сообщений из очереди
type Consumer interface {
	Consume(ctx context.Context, topicName string, handler consumer.Handler) (err error)
	Close() error
}
