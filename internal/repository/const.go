package repository

import (
	"errors"

	"github.com/redis/go-redis/v9"
)

var ErrFindByPhoneNotFound = errors.New("find by phone: row not found")

const RedisNil = redis.Nil
