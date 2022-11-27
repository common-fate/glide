package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/common-fate/pkg/providersetup"
	"github.com/common-fate/common-fate/pkg/storage/keys"
)

type ListProviderSetupSteps struct {
	SetupID string
	Result  []providersetup.Step `ddb:"result"`
}

func (l *ListProviderSetupSteps) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := dynamodb.QueryInput{
		KeyConditionExpression: aws.String("PK = :pk1 and begins_with(SK, :sk1)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk1": &types.AttributeValueMemberS{Value: keys.ProviderSetupStep.PK1},
			":sk1": &types.AttributeValueMemberS{Value: l.SetupID},
		},
	}
	return &qi, nil
}
