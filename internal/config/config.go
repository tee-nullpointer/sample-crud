package config

import (
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/tee-nullpointer/go-common-kit/pkg/env"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Logger   LoggerConfig
	Grpc     GrpcConfig
}

type ServerConfig struct {
	Host string
	Port string
	Mode string
}

type DatabaseConfig struct {
	Host          string
	Port          string
	Name          string
	User          string
	Password      string
	SSLMode       string
	MaxConnection int
	MaxIdle       int
	MaxLifetime   time.Duration
	MaxIdleTime   time.Duration
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
	PoolSize int
	MinIdle  int
	MaxWait  time.Duration
	IdleTime time.Duration
}

type LoggerConfig struct {
	Level  string
	Format string
}

type GrpcConfig struct {
	Port int
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	return &Config{
		Server: ServerConfig{
			Host: env.GetEnv("SERVER_HOST", "localhost"),
			Port: env.GetEnv("SERVER_PORT", "8080"),
			Mode: env.GetEnv("SERVER_MODE", "release"),
		},
		Database: DatabaseConfig{
			Host:          env.GetEnv("DATABASE_HOST", "localhost"),
			Port:          env.GetEnv("DATABASE_PORT", "5432"),
			Name:          env.GetEnv("DATABASE_NAME", "postgres"),
			User:          env.GetEnv("DATABASE_USER", "postgres"),
			Password:      env.GetEnv("DATABASE_PASSWORD", "postgres"),
			SSLMode:       env.GetEnv("DATABASE_SSL_MODE", "disable"),
			MaxConnection: env.GetEnvAsInt("DATABASE_MAX_CONNECTION", 10),
			MaxIdle:       env.GetEnvAsInt("DATABASE_MAX_IDLE", 5),
			MaxLifetime:   env.GetEnvAsDuration("DATABASE_MAX_LIFETIME", time.Minute*30),
			MaxIdleTime:   env.GetEnvAsDuration("DATABASE_MAX_IDLE_TIME", time.Minute*10),
		},
		Redis: RedisConfig{
			Host:     env.GetEnv("REDIS_HOST", "localhost"),
			Port:     env.GetEnv("REDIS_PORT", "6379"),
			Password: env.GetEnv("REDIS_PASSWORD", ""),
			DB:       env.GetEnvAsInt("REDIS_DB", 0),
			PoolSize: env.GetEnvAsInt("REDIS_POOL_SIZE", 20),
			MinIdle:  env.GetEnvAsInt("REDIS_MIN_IDLE", 5),
			MaxWait:  env.GetEnvAsDuration("REDIS_MAX_WAIT", time.Second*5),
			IdleTime: env.GetEnvAsDuration("REDIS_IDLE_TIME", time.Minute*30),
		},
		Logger: LoggerConfig{
			Level:  env.GetEnv("LOG_LEVEL", "info"),
			Format: env.GetEnv("LOG_FORMAT", "json"),
		},
		Grpc: GrpcConfig{
			Port: env.GetEnvAsInt("GRPC_PORT", 9090),
		},
	}
}
