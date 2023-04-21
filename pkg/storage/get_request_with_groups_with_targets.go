package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/storage/keys"
)

type GetRequestWithGroupsWithTargets struct {
	ID     string
	Result *access.RequestWithGroupsWithTargets
}

func (g *GetRequestWithGroupsWithTargets) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := &dynamodb.QueryInput{
		KeyConditionExpression: aws.String("PK = :pk1 and begins_with(SK, :sk1)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk1": &types.AttributeValueMemberS{Value: keys.AccessRequest.PK1},
			":sk1": &types.AttributeValueMemberS{Value: keys.AccessRequest.SK1(g.ID)},
		},
	}

	return qi, nil
}

func (g *GetRequestWithGroupsWithTargets) UnmarshalQueryOutput(out *dynamodb.QueryOutput) error {
	result, err := UnmarshalRequest(out.Items)
	if err != nil {
		return err
	}
	g.Result = result
	return nil
}
