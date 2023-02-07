package repository

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	log "github.com/sirupsen/logrus"
	"github.com/vaibhavahuja/rate-limiter/internal/app/entities"
	"strconv"
)

type RuleRepository Repository

func NewRuleRepository(client *dynamodb.Client, tableName string) *RuleRepository {
	return &RuleRepository{
		client:    client,
		tableName: tableName,
	}
}

// GetRuleByServiceId fetches the rate limiting rule by the service id
func (rr *RuleRepository) GetRuleByServiceId(ctx context.Context, serviceId int) (entities.Rule, error) {
	//todo this method needs some refining, will work on it
	//serviceIdItemInput, err := attributevalue.Marshal(serviceId)
	serviceIdItemInput := &types.AttributeValueMemberN{Value: strconv.Itoa(serviceId)}

	getItemInput := &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"service_id": serviceIdItemInput,
		},
		TableName: aws.String(rr.tableName),
	}

	log.Info("Fetching rule from table")
	response, err := rr.client.GetItem(ctx, getItemInput)
	if err != nil {
		log.Errorf("error while fetching item : %s", err.Error())
		return entities.Rule{}, err
	}
	if len(response.Item) == 0 {
		log.Infof("no values exist right now")
		return entities.Rule{}, nil

	}
	log.Infof("received item is %v", response.Item)

	var rule entities.ServiceEntity

	err = attributevalue.UnmarshalMap(response.Item, &rule)
	if err != nil {
		log.Errorf("unable to unmarshal response. here's why : %s", err.Error())
	}
	return rule.Rule, nil
}

// AddService saves the service information in ddb
func (rr *RuleRepository) AddService(ctx context.Context, serviceId int, rule entities.Rule) error {

	item, err := attributevalue.MarshalMap(entities.ServiceEntity{
		ServiceId: serviceId,
		Rule:      rule,
	})

	if err != nil {
		log.Errorf("error while saving rule %s", err.Error())
		return err
	}
	_, err = rr.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(rr.tableName),
		Item:      item,
	})
	if err != nil {
		log.Errorf("unable to add item to table. here's why : %v", err.Error())
	}
	return err
}
