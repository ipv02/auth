package redis

import (
	"context"
	"strconv"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/ipv02/auth/internal/client/cache"
	"github.com/ipv02/auth/internal/model"
	"github.com/ipv02/auth/internal/repository"
	"github.com/ipv02/auth/internal/repository/user/redis/converter"
	modelRepo "github.com/ipv02/auth/internal/repository/user/redis/model"
)

type repo struct {
	cl cache.RedisClient
}

// NewRepository создает новый экземпляр репозитория и возвращает его как интерфейс
func NewRepository(cl cache.RedisClient) repository.UserRepository {
	return &repo{cl: cl}
}

func (r *repo) CreateUser(ctx context.Context, user *model.UserCreate) (int64, error) {
	id := int64(1)

	userCreate := modelRepo.UserCreate{
		Name:            user.Name,
		Email:           user.Email,
		Password:        user.Password,
		PasswordConfirm: user.PasswordConfirm,
		Role:            user.Role,
	}

	idStr := strconv.FormatInt(id, 10)
	err := r.cl.HashSet(ctx, idStr, userCreate)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *repo) GetUser(ctx context.Context, id int64) (*model.UserGet, error) {
	idStr := strconv.FormatInt(id, 10)
	values, err := r.cl.HGetAll(ctx, idStr)
	if err != nil {
		return nil, err
	}

	if len(values) == 0 {
		return nil, model.ErrorUserNotFound
	}

	var user modelRepo.User
	err = redigo.ScanStruct(values, &user)
	if err != nil {
		return nil, err
	}

	return converter.ToUserFromRepo(&user), nil
}

func (r *repo) UpdateUser(ctx context.Context, user *model.UserUpdate) error {
	idStr := strconv.FormatInt(user.ID, 10)

	userUpdate := modelRepo.UserUpdate{
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}

	err := r.cl.HashSet(ctx, idStr, userUpdate)
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) DeleteUser(ctx context.Context, id int64) error {
	idStr := strconv.FormatInt(id, 10)
	return r.cl.Delete(ctx, idStr)
}
