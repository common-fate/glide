package storage

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/ddb"
	"github.com/common-fate/granted-approvals/pkg/rule"
	"github.com/common-fate/granted-approvals/pkg/storage/keys"
)

type ListAccessRulesForGroupsAndStatus struct {
	Groups []string
	Status rule.Status
	Result []rule.AccessRule `ddb:"result"`
}

func (l *ListAccessRulesForGroupsAndStatus) BuildQuery() (*dynamodb.QueryInput, error) {
	if len(l.Groups) == 0 {
		// return early with empty as we won't get any rules if no groups are provided.
		return nil, ddb.ErrNoItems
	}

	qi := dynamodb.QueryInput{
		IndexName:              &keys.IndexNames.GSI1,
		KeyConditionExpression: aws.String("GSI1PK = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: keys.AccessRule.GSI1PK(string(l.Status))},
		},
	}

	var expr string
	for i, g := range l.Groups {
		key := fmt.Sprintf(":group_%d", i)
		expr += fmt.Sprintf("contains(groups, %s)", key)
		if i < len(l.Groups)-1 {
			expr += " OR "
		}
		qi.ExpressionAttributeValues[key] = &types.AttributeValueMemberS{Value: g}
	}
	qi.FilterExpression = &expr
	return &qi, nil
}
