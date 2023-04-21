package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/storage/keys"
)

type GetRequestWithGroupsWithTargetsForUser struct {
	UserID    string
	RequestID string
	Result    *access.RequestWithGroupsWithTargets
}

func (g *GetRequestWithGroupsWithTargetsForUser) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := &dynamodb.QueryInput{
		IndexName:              &keys.IndexNames.GSI1,
		KeyConditionExpression: aws.String("GSI1PK = :pk1 and begins_with(GSI1SK, :sk1)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk1": &types.AttributeValueMemberS{Value: keys.AccessRequest.GSI1PK(g.UserID)},
			":sk1": &types.AttributeValueMemberS{Value: keys.AccessRequest.GSI1SKRequest(g.RequestID)},
		},
	}

	return qi, nil
}

func (g *GetRequestWithGroupsWithTargetsForUser) UnmarshalQueryOutput(out *dynamodb.QueryOutput) error {
	result, err := UnmarshalRequest(out.Items)
	if err != nil {
		return err
	}
	g.Result = result
	return nil
}
