package db

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)

type RedisDB struct {
	redisOptions *redis.Options
	redisClient  *redis.Client
}

func NewRedisDB(ctx context.Context, uri string) (*RedisDB, error) {
	redisOptions := &redis.Options{
		Addr: uri,
	}

	redisClient := redis.NewClient(redisOptions)

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	log.Debug().Msg("Redis connected!")
	
	return &RedisDB{
		redisOptions: redisOptions,
		redisClient:  redisClient,
	}, nil
}

func(redisDB *RedisDB) Close() error {
	return redisDB.redisClient.Close()
}