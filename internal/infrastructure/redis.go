package infrastructure

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedis(redisAddr string) (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	// проверить соединение
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return redisClient, nil
}
