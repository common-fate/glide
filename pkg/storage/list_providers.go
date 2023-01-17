package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	awsTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/common-fate/pkg/types"
)

type ListProviders struct {
	Result []types.ProviderV2 `ddb:"result"`
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
