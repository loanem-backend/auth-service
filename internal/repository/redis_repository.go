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
	Delete(ctx context.Context, key string) error
	StorePasswordChange(ctx context.Context, key string, value RedisPasswordChange, dur time.Duration) error
	GetPasswordChange(ctx context.Context, key string) (*RedisPasswordChange, error)
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

func (r *redisRepository) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *redisRepository) StorePasswordChange(ctx context.Context, key string, value RedisPasswordChange, dur time.Duration) error {
	if err := r.client.HSet(ctx, key, value).Err(); err != nil {
		return err
	}

	if err := r.client.Expire(ctx, key, dur).Err(); err != nil {
		return err
	}

	return nil
}

type RedisPasswordChange struct {
	HashedNewPassword string `redis:"hashed_new_password"`
	AssistantID       int    `redis:"assistant_id"`
}

func (r *redisRepository) GetPasswordChange(ctx context.Context, key string) (*RedisPasswordChange, error) {
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if exists < 1 {
		return nil, RedisNil
	}

	var value RedisPasswordChange

	if err := r.client.HGetAll(ctx, key).Scan(&value); err != nil {
		return nil, err
	}

	return &value, nil
}
