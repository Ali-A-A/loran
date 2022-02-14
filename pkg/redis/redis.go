package redis

import (
	"context"

	"github.com/Ali-A-A/loran/config"

	"github.com/go-redis/redis/v8"
)

// NewRedisClient returns new redis client
func NewRedisClient(redisConfig config.Redis) (*redis.Client, error) {
	redisClient := redis.NewClient(
		&redis.Options{
			Addr:            redisConfig.Master.Address,
			PoolSize:        redisConfig.Master.PoolSize,
			DialTimeout:     redisConfig.Master.DialTimeout,
			ReadTimeout:     redisConfig.Master.ReadTimeout,
			WriteTimeout:    redisConfig.Master.WriteTimeout,
			PoolTimeout:     redisConfig.Master.PoolTimeout,
			IdleTimeout:     redisConfig.Master.IdleTimeout,
			MinIdleConns:    redisConfig.Master.MinIdleConns,
			MaxRetries:      redisConfig.Master.MaxRetries,
			MinRetryBackoff: redisConfig.Master.MinRetryBackoff,
			MaxRetryBackoff: redisConfig.Master.MaxRetryBackoff,
		},
	)

	_, err := redisClient.Ping(context.TODO()).Result()
	if err != nil {
		return nil, err
	}

	return redisClient, nil
}
