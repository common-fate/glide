package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/common-fate/pkg/providersetup"
	"github.com/common-fate/common-fate/pkg/storage/keys"
)

type ListProviderSetupsForType struct {
	Type   string
	Result []providersetup.Setup `ddb:"result"`
}

func (l *ListProviderSetupsForType) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := dynamodb.QueryInput{
		IndexName:              aws.String(keys.IndexNames.GSI1),
		KeyConditionExpression: aws.String("GSI1PK = :pk1 and begins_with(GSI1SK, :sk1)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk1": &types.AttributeValueMemberS{Value: keys.ProviderSetup.GSI1PK},
			":sk1": &types.AttributeValueMemberS{Value: l.Type},
		},
	}
	return &qi, nil
}
