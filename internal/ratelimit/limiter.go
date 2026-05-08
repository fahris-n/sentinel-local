package ratelimit

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Limiter struct {
	Client *redis.Client
	Script *redis.Script
}

func NewLimiter(client *redis.Client, script *redis.Script) *Limiter {
	return &Limiter{
		Client: client,
		Script: script,
	}
}

func (l *Limiter) Allow(userId string, maxTokens int16, refillRate int16) (bool, error) {
	res, err := l.Script.Run(context.Background(), l.Client, []string{userId}, maxTokens, refillRate).Result()
	if err != nil {
		return false, fmt.Errorf("failed to query redis: %w", err)
	}

	if res.(int64) != 1 {
		return false, nil
	}

	return true, nil
}
