package repository

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRepository interface {
	Store(ctx context.Context, key string, value any, dur time.Duration) error
	GetInt(ctx context.Context, key string) (int, error)
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

func (r *redisRepository) GetInt(ctx context.Context, key string) (int, error) {
	val, err := r.client.Get(ctx, key).Int()

	if errors.Is(err, redis.Nil) {
		err = RedisNil
	}

	return val, err
}
