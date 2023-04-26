package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/ddb"
)

type GetRequestGroupWithTargetsForReviewer struct {
	RequestID  string
	GroupID    string
	ReviewerID string
	Result     *access.GroupWithTargets
}

func (g *GetRequestGroupWithTargetsForReviewer) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := &dynamodb.QueryInput{
		// Note that no limit(1) is used here because we need to fetch more that one items to read the whole request
		KeyConditionExpression: aws.String("PK = :pk1 and begins_with(SK, :sk1)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk1":        &types.AttributeValueMemberS{Value: keys.AccessRequestGroup.PK1},
			":sk1":        &types.AttributeValueMemberS{Value: keys.AccessRequestGroup.SK1(g.RequestID, g.GroupID)},
			":reviewerId": &types.AttributeValueMemberS{Value: g.ReviewerID},
		},
		FilterExpression: aws.String("contains(requestReviewers, :reviewerId)"),
	}

	return qi, nil
}

func (g *GetRequestGroupWithTargetsForReviewer) UnmarshalQueryOutput(out *dynamodb.QueryOutput) (*ddb.UnmarshalResult, error) {
	result, err := UnmarshalRequestGroup(out.Items)
	if err != nil {
		return nil, err
	}
	g.Result = result
	return &ddb.UnmarshalResult{}, nil
}
