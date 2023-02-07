package infrastructure

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	log "github.com/sirupsen/logrus"
	config2 "github.com/vaibhavahuja/rate-limiter/config"
)

func GetDynamoDBClient(region, url string) *dynamodb.Client {
	conf := config2.GetConfig()
	credProvider := aws.CredentialsProviderFunc(func(context.Context) (aws.Credentials, error) {
		return aws.Credentials{
			AccessKeyID:     conf.AWSConfig.AwsAccessKeyId,
			SecretAccessKey: conf.AWSConfig.AwsSecretAccessKey,
		}, nil
	})

	sdkConfig, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credProvider),
		config.WithEndpointResolver(aws.EndpointResolverFunc(
			func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{URL: url}, nil
			})),
	)
	if err != nil {
		log.Fatalf("Unable to load Dynamo DB Client. Here's why %s", err)
	}

	client := dynamodb.NewFromConfig(sdkConfig)
	log.Infof("Created dynamoDB client for region %s and url %s", region, url)
	return client
}
