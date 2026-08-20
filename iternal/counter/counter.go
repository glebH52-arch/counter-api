package counter

import (
	"context"
	"errors"
	"strconv"

	"github.com/redis/go-redis/v9"
)

const counterKey = "counter"

type CounterIntarface interface {
	IncrCount(ctx context.Context) (int64, error)
	GetCount(ctx context.Context) (int64, error)
}

type CounterRedisService struct {
	RedisClient *redis.Client
}

func NewCounterRedisService(redisClient *redis.Client) *CounterRedisService {
	return &CounterRedisService{
		RedisClient: redisClient,
	}
}

func (c *CounterRedisService) IncrCount(ctx context.Context) (int64, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}
	count, err := c.RedisClient.Incr(ctx, counterKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, nil
		}
		return 0, err
	}
	return count, nil
}

func (c *CounterRedisService) GetCount(ctx context.Context) (int64, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}
	valueStr, err := c.RedisClient.Get(ctx, counterKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, nil
		}
		return 0, err
	}
	value, _ := strconv.ParseInt(valueStr, 10, 64)
	return value, nil
}
