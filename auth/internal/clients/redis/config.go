package goredis

import (
	"fmt"
	"net"
	"os"
)

const (
	PortEnvName = "REDIS_PORT"
	HostEnvName = "REDIS_HOST"
)

type IRedisConfig interface {
	Address() string
}

type redisConfig struct {
	host string
	port string
}

func NewRedisConfig() (IRedisConfig, error) {
	port := os.Getenv(PortEnvName)
	host := os.Getenv(HostEnvName)
	if len(port) == 0 || len(host) == 0 {
		return nil, fmt.Errorf("REDIS_PORT or REDIS_HOST is not set")
	}

	return &redisConfig{
		host: host,
		port: port,
	}, nil
}

func (m *redisConfig) Address() string {
	return net.JoinHostPort(m.host, m.port)
}
