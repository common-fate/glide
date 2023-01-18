package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	awsTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/common-fate/pkg/provider"
	"github.com/common-fate/common-fate/pkg/storage/keys"
)

type ListProviders struct {
	Result []provider.Provider `ddb:"result"`
}

func (l *ListProviders) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := dynamodb.QueryInput{
		KeyConditionExpression: aws.String("PK = :pk1"),
		ExpressionAttributeValues: map[string]awsTypes.AttributeValue{
			":pk1": &awsTypes.AttributeValueMemberS{Value: keys.Provider.PK1},
		},
	}

	return &qi, nil
}
