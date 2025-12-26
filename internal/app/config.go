package app

import (
	"event-driven-architecture/pkg/envloader"
	"time"
)

type Config struct {
	LogLevel string

	HTTP     HTTPConfig
	SSE      SSEConfig
	App      AppConfig
	Security SecurityConfig

	Postgres PostgresConfig
	Redis    RedisConfig
}

type HTTPConfig struct {
	Port int
}

func NewHTTPConfig() HTTPConfig {
	return HTTPConfig{
		Port: envloader.MustGetInt("HTTP_PORT"),
	}
}

type SSEConfig struct {
	Port int
}

func NewSSEConfig() SSEConfig {
	return SSEConfig{
		Port: envloader.MustGetInt("SSE_PORT"),
	}
}

type AppConfig struct {
	ContextTimeout time.Duration
}

func NewAppConfig() AppConfig {
	return AppConfig{
		ContextTimeout: envloader.MustGetDuration("APP_CTX_TIMEOUT"),
	}
}

type SecurityConfig struct {
	BcryptCost int
}

func NewSecurityConfig() SecurityConfig {
	return SecurityConfig{
		BcryptCost: envloader.MustGetInt("BCRYPT_COST"),
	}
}

type PostgresConfig struct {
	DNS string
}

func NewPostgresConfig() PostgresConfig {
	return PostgresConfig{
		DNS: envloader.MustGetString("POSTGRES_DNS"),
	}
}

type RedisConfig struct {
	Addr string
}

func NewRedisConfig() RedisConfig {
	return RedisConfig{
		Addr: envloader.MustGetString("REDIS_ADDR"),
	}
}

func LoadConfig() Config {
	return Config{
		LogLevel: envloader.MustGetString("LOG_LEVEL"),
		HTTP:     NewHTTPConfig(),
		SSE:      NewSSEConfig(),
		App:      NewAppConfig(),
		Security: NewSecurityConfig(),
		Postgres: NewPostgresConfig(),
		Redis:    NewRedisConfig(),
	}
}
