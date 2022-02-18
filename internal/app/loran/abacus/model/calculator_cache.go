package model

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

// InMemoryCalculator is an in memory implementation of CalculatorRepo.
type InMemoryCalculator struct {
	rc *redis.Client
}

func getKey(entityID int64) string {
	return fmt.Sprintf("count:%d", entityID)
}

// NewInMemoryCalculator returns new InMemoryTrackCache
func NewInMemoryCalculator(rc *redis.Client) *InMemoryCalculator {
	return &InMemoryCalculator{
		rc: rc,
	}
}

// Add uses HyperLogLog algorithm which has been implemented in redis.
func (c *InMemoryCalculator) Add(ctx context.Context, userID int32, entityID int64) error {
	key := getKey(entityID)

	return c.rc.PFAdd(ctx, key, userID).Err()
}

// Count returns current counter of the entity.
func (c *InMemoryCalculator) Count(ctx context.Context, entityID int64) (int64, error) {
	key := getKey(entityID)

	cnt, err := c.rc.PFCount(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	return cnt, nil
}
