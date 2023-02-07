package config

import (
	"sync"
)

var (
	config *Config
	once   sync.Once
)

type Service struct {
	Name        string
	Port        string
	LogFileName string
	LogLevel    string
}

type AWSConfig struct {
	AwsAccessKeyId     string
	AwsSecretAccessKey string
}

type Dynamo struct {
	Region   string
	Endpoint string
	Table    string
	Index    string
}

type Config struct {
	Service           Service
	AWSConfig         AWSConfig
	RateLimiterDynamo Dynamo
}

func GetConfig() *Config {
	once.Do(initConfig)
	return config
}

func initConfig() {
	config = &Config{
		Service: Service{
			Name:        "rate-limiter",
			Port:        "8080",
			LogFileName: "./var/log/",
			LogLevel:    "debug",
		},
		AWSConfig: AWSConfig{
			AwsAccessKeyId:     "test",
			AwsSecretAccessKey: "temp",
		},
		RateLimiterDynamo: Dynamo{
			Region:   "temp",
			Endpoint: "temp",
			Table:    "temp",
			Index:    "temp",
		},
	}
}
