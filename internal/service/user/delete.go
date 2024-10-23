package user

import "context"

func (s *serv) DeleteUser(ctx context.Context, id int64) error {
	err := s.userRepository.DeleteUser(ctx, id)

	return err
}
