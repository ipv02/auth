package app

import (
	"context"
	"log"

	"github.com/IBM/sarama"
	redigo "github.com/gomodule/redigo/redis"

	"github.com/ipv02/auth/internal/api/user"
	"github.com/ipv02/auth/internal/client/cache"
	"github.com/ipv02/auth/internal/client/cache/redis"
	"github.com/ipv02/auth/internal/client/db"
	"github.com/ipv02/auth/internal/client/db/pg"
	"github.com/ipv02/auth/internal/client/db/transaction"
	"github.com/ipv02/auth/internal/client/kafka"
	kafkaConsumer "github.com/ipv02/auth/internal/client/kafka/consumer"
	"github.com/ipv02/auth/internal/closer"
	"github.com/ipv02/auth/internal/config"
	"github.com/ipv02/auth/internal/config/env"
	"github.com/ipv02/auth/internal/repository"
	userRepository "github.com/ipv02/auth/internal/repository/user/pg"
	userRepositoryRedis "github.com/ipv02/auth/internal/repository/user/redis"
	"github.com/ipv02/auth/internal/service"
	userSaverConsumer "github.com/ipv02/auth/internal/service/consumer/user_saver"
	userService "github.com/ipv02/auth/internal/service/user"
)

type serviceProvider struct {
	pgConfig            config.PGConfig
	grpcConfig          config.GRPCConfig
	httpConfig          config.HTTPConfig
	swaggerConfig       config.SwaggerConfig
	redisConfig         config.RedisConfig
	storageConfig       config.StorageConfig
	kafkaConsumerConfig config.KafkaConsumerConfig

	dbClient  db.Client
	txManager db.TxManager

	redisPool   *redigo.Pool
	redisClient cache.RedisClient

	userRepository repository.UserRepository

	userService service.UserService

	userImpl *user.Implementation

	userSaverConsumer service.ConsumerService

	consumer             kafka.Consumer
	consumerGroup        sarama.ConsumerGroup
	consumerGroupHandler *kafkaConsumer.GroupHandler
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

// HTTPConfig представляет конфигурацию для подключения к http серверу
func (s *serviceProvider) HTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		cfg, err := env.NewHTTPConfig()
		if err != nil {
			log.Fatalf("failed to get http config: %s", err.Error())
		}

		s.httpConfig = cfg
	}

	return s.httpConfig
}

// SwaggerConfig представляет конфигурацию для подключения к swgger серверу
func (s *serviceProvider) SwaggerConfig() config.SwaggerConfig {
	if s.swaggerConfig == nil {
		cfg, err := env.NewSwaggerConfig()
		if err != nil {
			log.Fatalf("failed to get swagger config: %s", err.Error())
		}

		s.swaggerConfig = cfg
	}

	return s.swaggerConfig
}

// RedisConfig представляет конфигурацию для подключения к redis
func (s *serviceProvider) RedisConfig() config.RedisConfig {
	if s.redisConfig == nil {
		cfg, err := env.NewRedisConfig()
		if err != nil {
			log.Fatalf("failed to get redis config: %s", err.Error())
		}

		s.redisConfig = cfg
	}

	return s.redisConfig
}

// StorageConfig представляет конфигурацию для подключения к storage
func (s *serviceProvider) StorageConfig() config.StorageConfig {
	if s.storageConfig == nil {
		cfg, err := env.NewStorageConfig()
		if err != nil {
			log.Fatalf("failed to get storage config: %s", err.Error())
		}

		s.storageConfig = cfg
	}

	return s.storageConfig
}

// KafkaConsumerConfig представляет конфигурацию для подключения к kafka
func (s *serviceProvider) KafkaConsumerConfig() config.KafkaConsumerConfig {
	if s.kafkaConsumerConfig == nil {
		cfg, err := env.NewKafkaConsumerConfig()
		if err != nil {
			log.Fatalf("failed to get kafka consumer config: %s", err.Error())
		}

		s.kafkaConsumerConfig = cfg
	}

	return s.kafkaConsumerConfig
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

// RedisPool возвращает экземпляр пула редиса
func (s *serviceProvider) RedisPool() *redigo.Pool {
	if s.redisPool == nil {
		s.redisPool = &redigo.Pool{
			MaxIdle:     s.RedisConfig().MaxIdle(),
			IdleTimeout: s.RedisConfig().IdleTimeout(),
			DialContext: func(ctx context.Context) (redigo.Conn, error) {
				return redigo.DialContext(ctx, "tcp", s.RedisConfig().Address())
			},
		}
	}

	return s.redisPool
}

// RedisClient возвращает экземпляр клиента редиса
func (s *serviceProvider) RedisClient() cache.RedisClient {
	if s.redisClient == nil {
		s.redisClient = redis.NewClient(s.RedisPool(), s.RedisConfig())
	}

	return s.redisClient
}

// UserRepository возвращает экземпляр репозитория
func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		if s.StorageConfig().Mode() == "redis" {
			s.userRepository = userRepositoryRedis.NewRepository(s.RedisClient())
		}

		if s.StorageConfig().Mode() == "pg" {
			s.userRepository = userRepository.NewRepository(s.DBClient(ctx))
		}
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

// UserSaverConsumer возвращает экземпляр consumerService
func (s *serviceProvider) UserSaverConsumer(ctx context.Context) service.ConsumerService {
	if s.userSaverConsumer == nil {
		s.userSaverConsumer = userSaverConsumer.NewService(
			s.UserRepository(ctx),
			s.Consumer(),
		)
	}

	return s.userSaverConsumer
}

// Consumer создает consumer
func (s *serviceProvider) Consumer() kafka.Consumer {
	if s.consumer == nil {
		s.consumer = kafkaConsumer.NewConsumer(
			s.ConsumerGroup(),
			s.ConsumerGroupHandler(),
		)
		closer.Add(s.consumer.Close)
	}

	return s.consumer
}

// ConsumerGroup создает consumerGroup
func (s *serviceProvider) ConsumerGroup() sarama.ConsumerGroup {
	if s.consumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			s.KafkaConsumerConfig().Brokers(),
			s.KafkaConsumerConfig().GroupID(),
			s.KafkaConsumerConfig().Config(),
		)
		if err != nil {
			log.Fatalf("failed to create consumer group: %v", err)
		}

		s.consumerGroup = consumerGroup
	}

	return s.consumerGroup
}

// ConsumerGroupHandler создает consumerGroupHandler
func (s *serviceProvider) ConsumerGroupHandler() *kafkaConsumer.GroupHandler {
	if s.consumerGroupHandler == nil {
		s.consumerGroupHandler = kafkaConsumer.NewGroupHandler()
	}

	return s.consumerGroupHandler
}
