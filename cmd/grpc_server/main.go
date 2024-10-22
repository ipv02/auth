package main

import (
	"context"
	"flag"
	"log"
	"net"

	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/ipv02/auth/internal/config"
	"github.com/ipv02/auth/internal/config/env"
	"github.com/ipv02/auth/internal/repository"
	"github.com/ipv02/auth/internal/repository/user"
	"github.com/ipv02/auth/pkg/user_v1"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
	user_v1.UnimplementedUserV1Server
	userRepository repository.UserRepository
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

	userRepo := user.NewRepository(pool)

	s := grpc.NewServer()
	reflection.Register(s)
	user_v1.RegisterUserV1Server(s, &server{userRepository: userRepo})

	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// CreateUser - запрос создает нового пользователя.
func (s *server) CreateUser(ctx context.Context, req *user_v1.CreateUserRequest) (*user_v1.CreateUserResponse, error) {
	userObj, err := s.userRepository.CreateUser(ctx, req)
	if err != nil {
		return nil, err
	}

	log.Printf("created user: %v", userObj)

	return &user_v1.CreateUserResponse{
		Id: userObj.Id,
	}, nil
}

// GetUser запроолс получения информации о пользователе.
func (s *server) GetUser(ctx context.Context, req *user_v1.GetUserRequest) (*user_v1.GetUserResponse, error) {
	userObj, err := s.userRepository.GetUser(ctx, req)
	if err != nil {
		return nil, err
	}

	log.Printf("get user: %v", userObj)

	return &user_v1.GetUserResponse{
		Id:        userObj.Id,
		Name:      userObj.Name,
		Email:     userObj.Email,
		Role:      userObj.Role,
		CreatedAt: userObj.CreatedAt,
		UpdatedAt: userObj.UpdatedAt,
	}, nil
}

// UpdateUser запрос на обновление данных о пользователе.
func (s *server) UpdateUser(ctx context.Context, req *user_v1.UpdateUserRequest) (*emptypb.Empty, error) {
	_, err := s.userRepository.UpdateUser(ctx, req)
	if err != nil {
		return nil, err
	}

	log.Printf("updated user: %v", req)

	return &emptypb.Empty{}, nil
}

// DeleteUser запрос на удаление пользователя.
func (s *server) DeleteUser(ctx context.Context, req *user_v1.DeleteUserRequest) (*emptypb.Empty, error) {
	_, err := s.userRepository.DeleteUser(ctx, req)
	if err != nil {
		return nil, err
	}

	log.Printf("deleted user: %v", req)

	return &emptypb.Empty{}, nil
}
