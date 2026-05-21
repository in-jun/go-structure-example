package database

import (
	"context"

	"github.com/in-jun/go-structure-example/internal/shared/config"

	"github.com/go-redis/redis/v8"
)

func NewRedis() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.AppConfig.RedisURL,
		DB:       0,
	})

	if _, err := client.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}

	return client, nil
}
