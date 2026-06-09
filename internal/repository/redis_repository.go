package repository

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRepository interface {
	Store(ctx context.Context, key string, value any, dur time.Duration) error
	GetString(ctx context.Context, key string) (string, error)
}

type redisRepository struct {
	client *redis.Client
}

func NewRedisRepository(rc *redis.Client) RedisRepository {
	return &redisRepository{
		client: rc,
	}
}

func (r *redisRepository) Store(ctx context.Context, key string, value any, dur time.Duration) error {
	return r.client.Set(ctx, key, value, dur).Err()
}

func (r *redisRepository) GetString(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()

	if errors.Is(err, redis.Nil) {
		err = RedisNil
	}

	return val, err
}

const RedisNil = redis.Nil
