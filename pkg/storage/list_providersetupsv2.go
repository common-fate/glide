package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/common-fate/pkg/providersetupv2"
	"github.com/common-fate/common-fate/pkg/storage/keys"
)

type ListProviderSetupsV2 struct {
	Result []providersetupv2.Setup `ddb:"result"`
}

func (l *ListProviderSetupsV2) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := dynamodb.QueryInput{
		KeyConditionExpression: aws.String("PK = :pk1"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk1": &types.AttributeValueMemberS{Value: keys.ProviderSetupV2.PK1},
		},
	}
	return &qi, nil
}
