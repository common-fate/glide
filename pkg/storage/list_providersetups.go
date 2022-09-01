package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/granted-approvals/pkg/providersetup"
	"github.com/common-fate/granted-approvals/pkg/storage/keys"
)

type ListProviderSetups struct {
	Result []providersetup.Setup `ddb:"result"`
}

func (l *ListProviderSetups) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := dynamodb.QueryInput{
		KeyConditionExpression: aws.String("PK = :pk1"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk1": &types.AttributeValueMemberS{Value: keys.ProviderSetup.PK1},
		},
	}
	return &qi, nil
}
