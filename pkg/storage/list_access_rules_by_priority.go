package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/common-fate/common-fate/pkg/storage/keys"
)

// ListAccessRulesByPriority lists rule in order from highest to lowest priority
type ListAccessRulesByPriority struct {
	Result []rule.AccessRule `ddb:"result"`
}

func (l *ListAccessRulesByPriority) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := dynamodb.QueryInput{
		IndexName:              &keys.IndexNames.GSI1,
		ScanIndexForward:       aws.Bool(false),
		KeyConditionExpression: aws.String("GSI1PK = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: keys.AccessRule.GSI1PK},
		},
	}
	return &qi, nil
}
