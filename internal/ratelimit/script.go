package ratelimit

import (
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

func LoadScript(filePath string) (*redis.Script, error) {
	luaFile, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read lua file: %w", err)
	}

	redisScript := redis.NewScript(string(luaFile))
	return redisScript, nil
}
