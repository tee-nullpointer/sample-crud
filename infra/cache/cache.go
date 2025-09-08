package cache

import (
	"context"
	"fmt"
	"sample-crud/internal/config"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var client *redis.Client

func NewRedisClient(cfg config.RedisConfig) *redis.Client {
	client = redis.NewClient(&redis.Options{
		Addr:            fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password:        cfg.Password,
		DB:              cfg.DB,
		PoolSize:        cfg.PoolSize,
		ConnMaxIdleTime: cfg.IdleTime,
		MinIdleConns:    cfg.MinIdle,
		PoolTimeout:     cfg.MaxWait,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Errorf("failed to connect to Redis: %w", err))
	}

	zap.L().Info("Successfully connected to Redis")
	return client
}

func Close() {
	err := client.Close()
	if err != nil {
		return
	}
}
