package service

import "github.com/vaibhavahuja/rate-limiter/config"

type Application struct {
	conf *config.Config
}

func NewApplication(conf *config.Config) *Application {
	return &Application{conf: conf}
}
