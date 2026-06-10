package config

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		value = defaultValue
	}

	return value
}

func InitDB() *pgxpool.Pool {
	ctx := context.Background()

	url := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		panic(fmt.Errorf("failed parsing databse config: %w", err))
	}
	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeExec

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		panic(fmt.Errorf("failed creating connection pool: %w", err))
	}

	return pool
}

func InitListener() net.Listener {
	listener, err := net.Listen("tcp", ":"+os.Getenv("APP_PORT"))
	if err != nil {
		panic(fmt.Errorf("failed listening: %w", err))
	}

	return listener
}

func InitRedisClient() *redis.Client {
	db_no, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		panic(fmt.Errorf("failed to connect to Redis: %w", err))
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASS"),
		DB:       db_no,
	})

	return rdb
}
