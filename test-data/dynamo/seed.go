package dynamo

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	log "github.com/sirupsen/logrus"
	"github.com/vaibhavahuja/rate-limiter/config"
	"github.com/vaibhavahuja/rate-limiter/internal/pkg/infrastructure"
)

func CreateServiceTable(conf *config.Config) (err error) {
	ctx := context.Background()
	ddbClient := infrastructure.GetDynamoDBClient(conf.RateLimiterDynamo.Region, conf.RateLimiterDynamo.Endpoint)
	_, err = ddbClient.CreateTable(ctx, getCreateServiceTableInput(conf.RateLimiterDynamo.Table))
	if err != nil {
		log.Errorf("error while creating table. here's why %s", err.Error())
		return
	}
	return
}

func getCreateServiceTableInput(tableName string) *dynamodb.CreateTableInput {
	return &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("service_id"),
				AttributeType: types.ScalarAttributeTypeN,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				KeyType:       types.KeyTypeHash,
				AttributeName: aws.String("service_id"),
			},
		},
		TableName: aws.String(tableName),
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
	}
}
