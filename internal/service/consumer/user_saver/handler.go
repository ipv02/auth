package user_saver

import (
	"context"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"

	"github.com/ipv02/auth/internal/model"
)

func (s *service) UserSaveHandler(ctx context.Context, msg *sarama.ConsumerMessage) error {
	userCreate := &model.UserCreate{}
	err := json.Unmarshal(msg.Value, userCreate)
	if err != nil {
		return err
	}

	id, err := s.userRepository.CreateUser(ctx, userCreate)
	if err != nil {
		return err
	}

	log.Printf("User with id %d created\n", id)

	return nil
}
