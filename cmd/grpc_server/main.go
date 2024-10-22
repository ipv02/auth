package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"net"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ipv02/auth/config"
	"github.com/ipv02/auth/config/env"
	"github.com/ipv02/auth/pkg/user_v1"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
	user_v1.UnimplementedUserV1Server
	pool *pgxpool.Pool
}

func main() {
	flag.Parse()
	ctx := context.Background()

	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	grpcConfig, err := env.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}

	pgConfig, err := env.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %v", err)
	}

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	pool, err := pgxpool.Connect(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer pool.Close()

	s := grpc.NewServer()
	reflection.Register(s)
	user_v1.RegisterUserV1Server(s, &server{pool: pool})

	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// CreateUser - запрос создает нового пользователя.
func (s *server) CreateUser(ctx context.Context, req *user_v1.CreateUserRequest) (*user_v1.CreateUserResponse, error) {
	log.Printf("CreateRequest - Name: %s, Email: %s", req.GetName(), req.GetEmail())

	builderInsert := sq.Insert("auth").
		PlaceholderFormat(sq.Dollar).
		Columns("name", "email", "password", "password_confirm", "role").
		Values(req.Name, req.Email, req.Password, req.PasswordConfirm, req.Role).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Fatalf("failed to generate query: %v", err)
		return nil, err
	}

	var userID int64
	err = s.pool.QueryRow(ctx, query, args...).Scan(&userID)

	if err != nil {
		log.Fatalf("failed to execute query: %v", err)
		return nil, err
	}

	return &user_v1.CreateUserResponse{
		Id: userID,
	}, nil
}

// GetUser запроолс получения информации о пользователе.
func (s *server) GetUser(ctx context.Context, req *user_v1.GetUserRequest) (*user_v1.GetUserResponse, error) {
	log.Printf("Auth id: %d", req.Id)

	builderSelect := sq.
		Select("id", "name", "email", "role", "created_at", "updated_at").
		From("auth").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": req.Id})

	query, args, err := builderSelect.ToSql()
	if err != nil {
		log.Fatalf("failed to generate query: %v", err)
		return nil, err
	}

	var (
		id          int64
		name, email string
		role        user_v1.UserRole
		createdAt   time.Time
		updatedAt   sql.NullTime
	)

	err = s.pool.QueryRow(ctx, query, args...).Scan(&id, &name, &email, &role, &createdAt, &updatedAt)
	if err != nil {
		log.Fatalf("failed to execute query: %v", err)
		return nil, err
	}

	var updatedAtProto *timestamppb.Timestamp
	if updatedAt.Valid {
		updatedAtProto = timestamppb.New(updatedAt.Time)
	}

	return &user_v1.GetUserResponse{
		Id:        id,
		Name:      name,
		Email:     email,
		Role:      role,
		CreatedAt: timestamppb.New(createdAt),
		UpdatedAt: updatedAtProto,
	}, nil
}

// UpdateUser запрос на обновление данных о пользователе.
func (s *server) UpdateUser(ctx context.Context, req *user_v1.UpdateUserRequest) (*emptypb.Empty, error) {
	log.Printf("UpdateRequest - ID: %d, Name: %s, Email: %s", req.GetId(), req.GetName(), req.GetEmail())

	builderUpdate := sq.
		Update("auth").
		PlaceholderFormat(sq.Dollar).
		Set("role", int(req.GetRole())).
		Set("updated_at", time.Now()).
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

	_, err = s.pool.Exec(ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to execute query: %v", err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// DeleteUser запрос на удаление пользователя.
func (s *server) DeleteUser(ctx context.Context, req *user_v1.DeleteUserRequest) (*emptypb.Empty, error) {
	log.Printf("Deleting object with ID: %d", req.GetId())

	builderDelete := sq.
		Delete("auth").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": req.Id})

	query, args, err := builderDelete.ToSql()
	if err != nil {
		log.Fatalf("failed to generate query: %v", err)
		return nil, err
	}

	_, err = s.pool.Exec(ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to execute query: %v", err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
