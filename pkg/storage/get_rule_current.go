package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/ddb"
	"github.com/common-fate/granted-approvals/pkg/rule"
	"github.com/common-fate/granted-approvals/pkg/storage/keys"
)

type GetAccessRuleCurrent struct {
	ID     string
	Result *rule.AccessRule
}

func (g *GetAccessRuleCurrent) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := &dynamodb.QueryInput{
		Limit:                  aws.Int32(1),
		KeyConditionExpression: aws.String("GSI2PK = :pk AND GSI2SK = :sk"),
		IndexName:              &keys.IndexNames.GSI2,
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: keys.AccessRule.GSI2PK},
			":sk": &types.AttributeValueMemberS{Value: keys.AccessRule.GSI2SK(g.ID)},
		},
	}
	return qi, nil
}

func (g *GetAccessRuleCurrent) UnmarshalQueryOutput(out *dynamodb.QueryOutput) error {
	if len(out.Items) != 1 {
		return ddb.ErrNoItems
	}

	return attributevalue.UnmarshalMap(out.Items[0], &g.Result)
}
