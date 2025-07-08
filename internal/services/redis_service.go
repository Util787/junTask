package service

import (
	"context"

	"github.com/Util787/user-manager-api/internal/repository"
)

type redisService struct {
	redisRepo repository.RedisRepository
}

func NewRedisService(repo repository.RedisRepository) RedisService {
	return &redisService{redisRepo: repo}
}

func (r *redisService) Set(ctx context.Context, key string, value any) error {
	return r.redisRepo.Set(ctx, key, value)
}


func (r *redisService) Get(ctx context.Context, key string, dest any) error {
	return r.redisRepo.Get(ctx, key, dest)
}

func (r *redisService) Delete(ctx context.Context, key string) error {
	return r.redisRepo.Delete(ctx, key)
}
