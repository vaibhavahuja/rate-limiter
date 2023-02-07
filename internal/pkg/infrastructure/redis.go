package infrastructure

import (
	"github.com/go-redis/redis"
)

func GetRedisClient(url, password string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: password,
		DB:       0,
	})
	return client
}
