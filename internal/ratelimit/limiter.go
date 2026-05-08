package ratelimit

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Limiter struct {
	client *redis.Client
	script *redis.Script
}

func NewLimiter(client *redis.Client, script *redis.Script) *Limiter {
	return &Limiter{
		client: client,
		script: script,
	}
}

func (l *Limiter) Allow(key string, maxTokens int16, refillRate int16) (bool, error) {
	res, err := l.script.Run(context.Background(), l.client, []string{key}, maxTokens, refillRate).Result()
	if err != nil {
		return false, fmt.Errorf("failed to query redis: %w", err)
	}

	if res.(int64) != 1 {
		return false, nil
	}

	return true, nil
}
