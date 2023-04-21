package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/ddb"
)

type ListRequestWithGroupsWithTargetsForUserAndPastUpcoming struct {
	UserID       string
	PastUpcoming keys.AccessRequestPastUpcoming
	Result       []access.RequestWithGroupsWithTargets
}

var _ ddb.QueryOutputUnmarshalerWithPagination = &ListRequestWithGroupsWithTargetsForUserAndPastUpcoming{}

func (g *ListRequestWithGroupsWithTargetsForUserAndPastUpcoming) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := &dynamodb.QueryInput{
		IndexName:              &keys.IndexNames.GSI2,
		KeyConditionExpression: aws.String("GSI2PK = :pk1"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk1": &types.AttributeValueMemberS{Value: keys.AccessRequest.GSI2PK(g.UserID, g.PastUpcoming)},
		},
	}

	return qi, nil
}

func (g *ListRequestWithGroupsWithTargetsForUserAndPastUpcoming) UnmarshalQueryOutputWithPagination(out *dynamodb.QueryOutput) (map[string]types.AttributeValue, error) {
	result, pagination, err := UnmarshalRequests(out.Items)
	if err != nil {
		return nil, err
	}
	g.Result = result
	return pagination, nil
}
