package queries

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/common-fate/cmd/gdeploy/commands/migrate/keys"
	"github.com/common-fate/common-fate/cmd/gdeploy/commands/migrate/rule"
)

type ListAccessRulesForStatus struct {
	Status rule.Status
	Result []rule.AccessRule `ddb:"result"`
}

func (l *ListAccessRulesForStatus) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := dynamodb.QueryInput{
		IndexName:              &keys.IndexNames.GSI1,
		KeyConditionExpression: aws.String("GSI1PK = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: keys.AccessRule.GSI1PK(string(l.Status))},
		},
	}
	return &qi, nil
}
