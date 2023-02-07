package service

import (
	"github.com/go-redis/redis"
	"github.com/vaibhavahuja/rate-limiter/config"
	"github.com/vaibhavahuja/rate-limiter/internal/app/gateway/repository"
)

type Application struct {
	conf *config.Config
	//will remove redis client from here as well
	redisClient     *redis.Client
	rulesRepository *repository.RuleRepository
}

func NewApplication(conf *config.Config, redisClient *redis.Client, rulesRepository *repository.RuleRepository) *Application {
	return &Application{
		conf:            conf,
		redisClient:     redisClient,
		rulesRepository: rulesRepository,
	}
}
