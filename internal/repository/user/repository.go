package user

import (
	"context"
	"log"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/ipv02/auth/internal/client/db"
	"github.com/ipv02/auth/internal/model"
	"github.com/ipv02/auth/internal/repository"
	"github.com/ipv02/auth/internal/repository/user/converter"
	modelRepo "github.com/ipv02/auth/internal/repository/user/model"
)

const (
	tableName = "auth"

	idColumn              = "id"
	nameColumn            = "name"
	emailColumn           = "email"
	passwordColumn        = "password"
	passwordConfirmColumn = "password_confirm"
	roleColumn            = "role"
	createdAtColumn       = "created_at"
	updatedAtColumn       = "updated_at"
)

type repo struct {
	db db.Client
}

// NewRepository создает новый экземпляр UserRepository с подключением к базе данных
func NewRepository(db db.Client) repository.UserRepository {
	return &repo{db: db}
}

// CreateUser выполняет создание нового пользователя в базе данных
func (r *repo) CreateUser(ctx context.Context, user *model.UserCreate) (int64, error) {
	builderInsert := sq.Insert(tableName).
		Columns(nameColumn, emailColumn, passwordColumn, passwordConfirmColumn, roleColumn).
		Values(user.Name, user.Email, user.Password, user.PasswordConfirm, user.Role).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Fatalf("failed to generate query: %v", err)
		return 0, err
	}

	q := db.Query{
		Name:     "user_repository.Create",
		QueryRaw: query,
	}

	var userID int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&userID)

	if err != nil {
		log.Fatalf("failed to execute query: %v", err)
		return 0, err
	}

	return userID, nil
}

// GetUser функция берет пользоваиеля из базы данных
func (r *repo) GetUser(ctx context.Context, id int64) (*model.UserGet, error) {
	builderSelect := sq.
		Select(idColumn, nameColumn, emailColumn, roleColumn, createdAtColumn, updatedAtColumn).
		From(tableName).
		Where(sq.Eq{idColumn: id}).
		PlaceholderFormat(sq.Dollar).
		Limit(1)

	query, args, err := builderSelect.ToSql()
	if err != nil {
		log.Fatalf("failed to generate query: %v", err)
		return nil, err
	}

	q := db.Query{
		Name:     "user_repository.Get",
		QueryRaw: query,
	}

	var user modelRepo.User
	err = r.db.DB().ScanOneContext(ctx, &user, q, args...)
	if err != nil {
		log.Fatalf("failed to execute query: %v", err)
		return nil, err
	}

	return converter.ToUserFromRepo(&user), nil
}

// UpdateUser выполняет обновление данных пользователя в базе
func (r *repo) UpdateUser(ctx context.Context, user *model.UserUpdate) error {
	builderUpdate := sq.
		Update(tableName).
		Set(roleColumn, user.Role).
		Set(updatedAtColumn, time.Now()).
		Where(sq.Eq{idColumn: user.ID}).
		PlaceholderFormat(sq.Dollar)

	if user.Name != nil {
		trimmedName := strings.TrimSpace(*user.Name)
		builderUpdate.Set(nameColumn, trimmedName)
	}

	if user.Email != nil {
		trimmedEmail := strings.TrimSpace(*user.Email)
		builderUpdate.Set(emailColumn, trimmedEmail)
	}

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		log.Fatalf("failed to generate query: %v", err)
		return err
	}

	q := db.Query{
		Name:     "user_repository.Update",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		log.Fatalf("failed to execute query: %v", err)
		return err
	}

	return err
}

// DeleteUser удаляет пользователя из базы данных
func (r *repo) DeleteUser(ctx context.Context, id int64) error {
	builderDelete := sq.
		Delete(tableName).
		Where(sq.Eq{idColumn: id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builderDelete.ToSql()
	if err != nil {
		log.Fatalf("failed to generate query: %v", err)
		return err
	}

	q := db.Query{
		Name:     "user_repository.Delete",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		log.Fatalf("failed to execute query: %v", err)
		return err
	}

	return err
}
