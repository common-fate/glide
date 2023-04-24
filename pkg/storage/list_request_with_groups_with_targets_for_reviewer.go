package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/ddb"
)

type ListRequestWithGroupsWithTargetsForReviewer struct {
	ReviewerID string
	Result     []access.RequestWithGroupsWithTargets
}

var _ ddb.QueryOutputUnmarshaler = &ListRequestWithGroupsWithTargetsForReviewer{}

func (g *ListRequestWithGroupsWithTargetsForReviewer) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := &dynamodb.QueryInput{
		ScanIndexForward:       aws.Bool(false),
		KeyConditionExpression: aws.String("PK = :pk1"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk1":        &types.AttributeValueMemberS{Value: keys.AccessRequest.PK1},
			":reviewerId": &types.AttributeValueMemberS{Value: g.ReviewerID},
		},
		FilterExpression: aws.String("contains(requestReviewers, :reviewerId)"),
	}

	return qi, nil
}

func (g *ListRequestWithGroupsWithTargetsForReviewer) UnmarshalQueryOutput(out *dynamodb.QueryOutput) (*ddb.UnmarshalResult, error) {
	result, pagination, err := UnmarshalRequestsBottomToTop(out.Items)
	if err != nil {
		return nil, err
	}
	g.Result = result
	return &ddb.UnmarshalResult{PaginationToken: pagination}, nil
}
