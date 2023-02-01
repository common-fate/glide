package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/common-fate/pkg/providerdeployment"
	"github.com/common-fate/common-fate/pkg/storage/keys"
)

type ListProviderDeployments struct {
	Result []providerdeployment.ProviderDeployment `ddb:"result"`
}

func (l *ListProviderDeployments) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := dynamodb.QueryInput{
		KeyConditionExpression: aws.String("PK = :pk1"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk1": &types.AttributeValueMemberS{Value: keys.ProviderDeployment.PK1},
		},
	}
	return &qi, nil
}
