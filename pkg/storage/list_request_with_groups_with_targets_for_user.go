package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/ddb"
)

type ListRequestWithGroupsWithTargetsForUser struct {
	UserID string
	Result []access.RequestWithGroupsWithTargets
}

func (g *ListRequestWithGroupsWithTargetsForUser) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := &dynamodb.QueryInput{
		ScanIndexForward:       aws.Bool(false),
		IndexName:              &keys.IndexNames.GSI1,
		KeyConditionExpression: aws.String("GSI1PK = :pk1"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk1": &types.AttributeValueMemberS{Value: keys.AccessRequest.GSI1PK(g.UserID)},
		},
	}

	return qi, nil
}

func (g *ListRequestWithGroupsWithTargetsForUser) UnmarshalQueryOutput(out *dynamodb.QueryOutput) (*ddb.UnmarshalResult, error) {
	result, pagination, err := UnmarshalRequestsBottomToTop(out.Items)
	if err != nil {
		return nil, err
	}
	g.Result = result
	return &ddb.UnmarshalResult{PaginationToken: pagination}, nil
}
