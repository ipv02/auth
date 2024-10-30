package user

import "context"

// DeleteUser запрос сервесного слоя на удаление пользователя
func (s *service) DeleteUser(ctx context.Context, id int64) error {
	err := s.userRepository.DeleteUser(ctx, id)

	return err
}
