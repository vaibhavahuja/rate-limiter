package cache

import (
	"github.com/go-redis/redis"
)

type Cache struct {
	redisClient *redis.Client
}
