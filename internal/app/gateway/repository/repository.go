package repository

import "github.com/aws/aws-sdk-go-v2/service/dynamodb"

type Repository struct {
	client    *dynamodb.Client
	tableName string
}
