package goredis

import (
	"context"
	"time"
)

type IRedis interface {
	Set(ctx context.Context, key string, value string, duration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	SetObject(ctx context.Context, key string, value interface{}, duration time.Duration) error
	GetObject(ctx context.Context, key string, value any) error
	Delete(ctx context.Context, key string) error
}
