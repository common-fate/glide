package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/common-fate/pkg/requestsv2.go"
	"github.com/common-fate/common-fate/pkg/storage/keys"
)

type ListEntitlementResources struct {
	Provider requestsv2.TargetFrom
	Argument string
	Result   []requestsv2.Option `ddb:"result"`
}

func (l *ListEntitlementResources) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := dynamodb.QueryInput{
		KeyConditionExpression: aws.String("PK = :pk1 and begins_with(SK, :sk1)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk1": &types.AttributeValueMemberS{Value: keys.OptionsV2.PK1},
			":sk1": &types.AttributeValueMemberS{Value: keys.OptionsV2.SK1All(l.Provider.GetTargetFromString())},
		},
	}
	return &qi, nil
}
