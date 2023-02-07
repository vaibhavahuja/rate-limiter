package config

import (
	"fmt"
	"os"
	"sync"
)

const (
	ServiceName = "SERVICE_NAME"
	GRPCServerPort = "SERVER_PORT"
	LogFileName = "LOG_FILE_NAME"
	LogLevel = "LOG_LEVEL"
	AwsAccessKeyId     = "AWS_ACCESS_KEY_ID"
	AwsSecretAccessKey = "AWS_SECRET_ACCESS_KEY"
	RateLimiterTableName     = "RATELIMITER_DB_TABLENAME"
	RateLimiterIndexName     = "RATELIMITER_DB_GSINAME"
	RateLimiterTableRegion   = "RATELIMITER_DB_REGION"
	RateLimiterTableEndpoint = "RATELIMITER_DB_ENDPOINT"
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
		Service:           Service{
			Name:        Get(ServiceName),
			Port:        Get(GRPCServerPort),
			LogFileName:        Get(LogFileName),
			LogLevel:           Get(LogLevel),
		},
		AWSConfig:         AWSConfig{
			AwsAccessKeyId:     Get(AwsAccessKeyId),
			AwsSecretAccessKey: Get(AwsSecretAccessKey),
		},
		RateLimiterDynamo: Dynamo{
			Region:   Get(RateLimiterTableRegion),
			Endpoint: Get(RateLimiterTableEndpoint),
			Table:    Get(RateLimiterTableName),
			Index:    Get(RateLimiterIndexName),
		},
	}
}


func Get(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	panic(fmt.Sprintf("Env Variable: %s is not defined", key))
}
