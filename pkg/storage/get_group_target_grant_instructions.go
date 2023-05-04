package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/ddb"
)

type GetGroupTargetGrantInstructions struct {
	UserId   string
	TargetID string
	Result   *access.Instructions
}

func (g *GetGroupTargetGrantInstructions) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := &dynamodb.QueryInput{

		KeyConditionExpression: aws.String("PK = :pk1 and begins_with(SK, :sk1)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk1": &types.AttributeValueMemberS{Value: keys.AccessRequestGroupTargetInstructions.PK1},
			":sk1": &types.AttributeValueMemberS{Value: keys.AccessRequestGroupTargetInstructions.SK1(g.TargetID, g.UserId)},
		},
	}

	return qi, nil
}

func (g *GetGroupTargetGrantInstructions) UnmarshalQueryOutput(out *dynamodb.QueryOutput) (*ddb.UnmarshalResult, error) {
	if len(out.Items) != 1 {
		return nil, ddb.ErrNoItems
	}

	return &ddb.UnmarshalResult{}, attributevalue.UnmarshalMap(out.Items[0], &g.Result)
}
