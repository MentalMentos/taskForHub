package goredis

import (
	"context"
	"encoding/json"
	"fmt"
	Redis "github.com/redis/go-redis/v9"
	"time"
)

type RedisClient struct {
	Client *Redis.Client
}

func NewRedisClient(config *redisConfig) *RedisClient {
	redisClient := &RedisClient{
		Client: Redis.NewClient(&Redis.Options{
			Addr:     config.Address(),
			Password: "",
			DB:       0,
		}),
	}

	return redisClient
}

func (redisClient *RedisClient) Set(ctx context.Context, key string, value interface{}) error {
	return redisClient.Client.Set(ctx, key, value, 0).Err()
}

func (redisClient *RedisClient) Get(ctx context.Context, key string) (interface{}, error) {
	return redisClient.Client.Get(ctx, key).Result()
}

func (redisClient *RedisClient) Del(ctx context.Context, key string) error {
	return redisClient.Client.Del(ctx, key).Err()
}

func (redisClient *RedisClient) SetObject(ctx context.Context, key string, value interface{}, duration time.Duration) error {
	noteBytes, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshall value in redis: %w", err)
	}
	redisClient.Client.Set(ctx, key, string(noteBytes), duration)
	return nil
}

func (redisClient *RedisClient) GetObject(ctx context.Context, key string, value interface{}) (interface{}, error) {
	val, err := redisClient.Client.Get(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get value from redis: %w", err)
	}
	if err = json.Unmarshal([]byte(val), &value); err != nil {
		return nil, fmt.Errorf("failed to unmarshall value in redis: %w", err)
	}
	return value, nil
}
