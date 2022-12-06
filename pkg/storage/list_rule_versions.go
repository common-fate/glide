package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/common-fate/common-fate/pkg/storage/keys"
)

type ListAccessRuleVersions struct {
	ID     string
	Result []rule.AccessRule `ddb:"result"`
}

func (l *ListAccessRuleVersions) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := dynamodb.QueryInput{
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :sk)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: keys.AccessRule.PK1},
			":sk": &types.AttributeValueMemberS{Value: keys.AccessRule.SK1RuleID(l.ID)},
		},
	}
	return &qi, nil
}
