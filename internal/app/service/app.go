package service

import (
	"github.com/vaibhavahuja/rate-limiter/config"
	"github.com/vaibhavahuja/rate-limiter/internal/app/gateway/cache"
	"github.com/vaibhavahuja/rate-limiter/internal/app/gateway/repository"
)

type Application struct {
	conf                        *config.Config
	rulesRepository             *repository.RuleRepository
	requestCounterCache         *cache.RequestCounterCache
}

func NewApplication(conf *config.Config, rulesRepository *repository.RuleRepository, requestCounterCache *cache.RequestCounterCache) *Application {
	return &Application{
		conf:                        conf,
		rulesRepository:             rulesRepository,
		requestCounterCache:         requestCounterCache,
	}
}
