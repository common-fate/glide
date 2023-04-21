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

var _ ddb.QueryOutputUnmarshalerWithPagination = &ListRequestWithGroupsWithTargetsForUser{}

func (g *ListRequestWithGroupsWithTargetsForUser) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := &dynamodb.QueryInput{
		IndexName:              &keys.IndexNames.GSI1,
		KeyConditionExpression: aws.String("GSI1PK = :pk1"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk1": &types.AttributeValueMemberS{Value: keys.AccessRequest.GSI1PK(g.UserID)},
		},
	}

	return qi, nil
}

func (g *ListRequestWithGroupsWithTargetsForUser) UnmarshalQueryOutputWithPagination(out *dynamodb.QueryOutput) (map[string]types.AttributeValue, error) {
	result, pagination, err := UnmarshalRequests(out.Items)
	if err != nil {
		return nil, err
	}
	g.Result = result
	return pagination, nil
}
