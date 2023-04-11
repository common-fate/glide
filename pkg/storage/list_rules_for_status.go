package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/common-fate/common-fate/pkg/storage/keys"
)

type ListAccessRulesForStatus struct {
	Status rule.Status
	Result []rule.AccessRule `ddb:"result"`
}

func (l *ListAccessRulesForStatus) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := dynamodb.QueryInput{
		IndexName:              &keys.IndexNames.GSI1,
		KeyConditionExpression: aws.String("GSI1PK = :pk AND GSI1SK = :sk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: keys.AccessRule.GSI1PK},
			":sk": &types.AttributeValueMemberS{Value: keys.AccessRule.GSI1SK(string(l.Status))},
		},
	}
	return &qi, nil
}
