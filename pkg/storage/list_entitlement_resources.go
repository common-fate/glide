package storage

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/common-fate/pkg/requestsv2.go"
	"github.com/common-fate/common-fate/pkg/storage/keys"
)

type ListEntitlementResources struct {
	Provider     requestsv2.TargetFrom
	Argument     string
	FilterValues []string
	Groups       []string

	Result []requestsv2.ResourceOption `ddb:"result"`
}

func (l *ListEntitlementResources) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := dynamodb.QueryInput{
		KeyConditionExpression: aws.String("PK = :pk1 and begins_with(SK, :sk1)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk1": &types.AttributeValueMemberS{Value: keys.OptionsV2.PK1(l.Argument)},
			":sk1": &types.AttributeValueMemberS{Value: keys.OptionsV2.SK1All(l.Provider.GetTargetFromString())},
		},
	}

	expr := "("
	for i, g := range l.Groups {
		key := fmt.Sprintf(":group_%d", i)
		expr += fmt.Sprintf("contains(accessRules, %s)", key)
		if i < len(l.Groups)-1 {
			expr += " OR "
		}
		qi.ExpressionAttributeValues[key] = &types.AttributeValueMemberS{Value: g}
	}
	expr += ")"

	//then we do an AND expression to filter for related fields if they exist

	if len(l.FilterValues) > 0 {
		expr += " AND ("
		for i, g := range l.FilterValues {
			key := fmt.Sprintf(":filter_%d", i)
			expr += fmt.Sprintf("contains(childOf, %s)", key)
			if i < len(l.FilterValues)-1 {
				expr += " OR "
			}
			qi.ExpressionAttributeValues[key] = &types.AttributeValueMemberS{Value: g}
		}
		expr += ")"
	}

	qi.FilterExpression = &expr
	return &qi, nil
}
