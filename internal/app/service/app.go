package service

import (
	"github.com/go-redis/redis"
	"github.com/vaibhavahuja/rate-limiter/config"
)

type Application struct {
	conf        *config.Config
	redisClient *redis.Client
}

func NewApplication(conf *config.Config, redisClient *redis.Client) *Application {
	return &Application{
		conf:        conf,
		redisClient: redisClient,
	}
}
