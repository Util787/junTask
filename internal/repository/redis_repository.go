package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

const redisTTL = 15 * time.Second

type redisRepository struct {
	redis *redis.Client
}

func NewRedisRepository(redis *redis.Client) RedisRepository {
	return &redisRepository{redis: redis}
}

func (r *redisRepository) Set(ctx context.Context, key string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.redis.Set(ctx, key, data, redisTTL).Err()
}

func (r *redisRepository) Get(ctx context.Context, key string, dest any) error {
	str, err := r.redis.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(str), dest)
}

func (r *redisRepository) Delete(ctx context.Context, key string) error {
	return r.redis.Del(ctx, key).Err()
}
