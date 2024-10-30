package app

import (
	"context"
	"log"

	"github.com/ipv02/auth/internal/api/user"
	"github.com/ipv02/auth/internal/client/db"
	"github.com/ipv02/auth/internal/client/db/pg"
	"github.com/ipv02/auth/internal/client/db/transaction"
	"github.com/ipv02/auth/internal/closer"
	"github.com/ipv02/auth/internal/config"
	"github.com/ipv02/auth/internal/config/env"
	"github.com/ipv02/auth/internal/repository"
	userRepository "github.com/ipv02/auth/internal/repository/user"
	"github.com/ipv02/auth/internal/service"
	userService "github.com/ipv02/auth/internal/service/user"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig

	dbClient       db.Client
	txManager      db.TxManager
	userRepository repository.UserRepository

	userService service.UserService

	userImpl *user.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

// PgConfig представляет конфигурацию для подключения к базе данных
func (s *serviceProvider) PgConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := env.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

// GRPCConfig представляет конфигурацию для подключения к gRPC серверу
func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := env.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

// DBClient клиент для работы с базой данных
func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PgConfig().DSN())
		if err != nil {
			log.Fatalf("failed to connect to database: %s", err.Error())
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("failed to ping database: %s", err.Error())
		}

		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

// TxManager возвращает экземпляр менеджера транзакций
func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

// UserRepository возвращает экземпляр репозитория
func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepository.NewRepository(s.DBClient(ctx))
	}

	return s.userRepository
}

// UserService возвращает экземпляр сервиса
func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userService.NewService(s.UserRepository(ctx), s.TxManager(ctx))
	}

	return s.userService
}

// UserImpl возвращает экземпляр имплементации
func (s *serviceProvider) UserImpl(ctx context.Context) *user.Implementation {
	if s.userImpl == nil {
		s.userImpl = user.NewImplementation(s.UserService(ctx))
	}

	return s.userImpl
}
