package config

import (
	"time"

	"github.com/joho/godotenv"
)

// Load загружает переменные окружения из указанного файла.
func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	return nil
}

// GRPCConfig представляет конфигурацию для подключения к gRPC серверу.
type GRPCConfig interface {
	Address() string
}

// HTTPConfig представляет конфигурацию для подключения к http серверу.
type HTTPConfig interface {
	Address() string
}

// SwaggerConfig представляет конфигурацию для подключения к swagger серверу.
type SwaggerConfig interface {
	Address() string
}

// PGConfig представляет конфигурацию для подключения к базе данных PostgreSQL.
type PGConfig interface {
	DSN() string
}

// RedisConfig представляет конфигурацию для подключения к redis
type RedisConfig interface {
	Address() string
	ConnectionTimeout() time.Duration
	MaxIdle() int
	IdleTimeout() time.Duration
}

// StorageConfig предназначен для конфигурации хранилища
type StorageConfig interface {
	Mode() string
}
