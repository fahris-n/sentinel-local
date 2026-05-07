package ratelimit

import (
	"github.com/redis/go-redis/v9"
)

type Limiter struct {
	Client *redis.Client
	Script *redis.Script
}
