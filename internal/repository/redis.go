package repository

import (
	"context"

	"github.com/Util787/user-manager-api/config"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(cfg config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	
	err := client.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}
	return client, nil
}
