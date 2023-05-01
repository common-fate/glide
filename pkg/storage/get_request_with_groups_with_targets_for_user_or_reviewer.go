package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/ddb"
)

type GetRequestWithGroupsWithTargetsForUserOrReviewer struct {
	// can be a user id or a reviewer id
	UserID    string
	RequestID string
	Result    *access.RequestWithGroupsWithTargets
}

func (g *GetRequestWithGroupsWithTargetsForUserOrReviewer) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := &dynamodb.QueryInput{
		// Note that no limit(1) is used here because we need to fetch more that one items to read the whole request
		KeyConditionExpression: aws.String("PK = :pk1 and begins_with(SK, :sk1)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk1":    &types.AttributeValueMemberS{Value: keys.AccessRequest.PK1},
			":sk1":    &types.AttributeValueMemberS{Value: keys.AccessRequest.SK1(g.RequestID)},
			":userId": &types.AttributeValueMemberS{Value: g.UserID},
		},
		FilterExpression: aws.String("requestedBy.id = :userId or contains(requestReviewers, :userId)"),
	}

	return qi, nil
}

func (g *GetRequestWithGroupsWithTargetsForUserOrReviewer) UnmarshalQueryOutput(out *dynamodb.QueryOutput) (*ddb.UnmarshalResult, error) {
	result, err := UnmarshalRequest(out.Items)
	if err != nil {
		return nil, err
	}
	g.Result = result
	return &ddb.UnmarshalResult{}, nil
}
