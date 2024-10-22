package user

import (
	"context"
	"log"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/ipv02/auth/internal/repository"
	"github.com/ipv02/auth/internal/repository/user/converter"
	"github.com/ipv02/auth/internal/repository/user/model"
	"github.com/ipv02/auth/pkg/user_v1"
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
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) repository.UserRepository {
	return &repo{db: db}
}

func (r *repo) CreateUser(ctx context.Context, req *user_v1.CreateUserRequest) (*user_v1.CreateUserResponse, error) {
	builderInsert := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(nameColumn, emailColumn, passwordColumn, passwordConfirmColumn, roleColumn).
		Values(req.Name, req.Email, req.Password, req.PasswordConfirm, req.Role).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Fatalf("failed to generate query: %v", err)
		return nil, err
	}

	var userID int64
	err = r.db.QueryRow(ctx, query, args...).Scan(&userID)

	if err != nil {
		log.Fatalf("failed to execute query: %v", err)
		return nil, err
	}

	return &user_v1.CreateUserResponse{
		Id: userID,
	}, nil
}

func (r *repo) GetUser(ctx context.Context, req *user_v1.GetUserRequest) (*user_v1.GetUserResponse, error) {
	builderSelect := sq.
		Select(idColumn, nameColumn, emailColumn, roleColumn, createdAtColumn, updatedAtColumn).
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": req.Id})

	query, args, err := builderSelect.ToSql()
	if err != nil {
		log.Fatalf("failed to generate query: %v", err)
		return nil, err
	}

	var user model.User
	err = r.db.QueryRow(ctx, query, args...).Scan(&user.ID, &user.Name, &user.Email, &user.UserRole, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		log.Fatalf("failed to execute query: %v", err)
		return nil, err
	}

	return converter.ToUserFromRepo(&user), nil
}

func (r *repo) UpdateUser(ctx context.Context, req *user_v1.UpdateUserRequest) (*emptypb.Empty, error) {
	builderUpdate := sq.
		Update(tableName).
		PlaceholderFormat(sq.Dollar).
		Set(roleColumn, int(req.GetRole())).
		Set(updatedAtColumn, time.Now()).
		Where(sq.Eq{"id": req.Id})

	if req.Name != nil {
		trimmedName := strings.TrimSpace(req.Name.GetValue())
		if len(trimmedName) > 0 {
			builderUpdate.Set("name", trimmedName)
		}
	}

	if req.Email != nil {
		trimmedEmail := strings.TrimSpace(req.Email.GetValue())
		if len(trimmedEmail) > 0 {
			builderUpdate.Set("email", trimmedEmail)
		}
	}

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		log.Fatalf("failed to generate query: %v", err)
		return nil, err
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to execute query: %v", err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (r *repo) DeleteUser(ctx context.Context, req *user_v1.DeleteUserRequest) (*emptypb.Empty, error) {
	builderDelete := sq.
		Delete(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": req.Id})

	query, args, err := builderDelete.ToSql()
	if err != nil {
		log.Fatalf("failed to generate query: %v", err)
		return nil, err
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to execute query: %v", err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
