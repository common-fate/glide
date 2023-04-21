package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/ddb"
)

type ListRequestWithGroupsWithTargets struct {
	Result []access.RequestWithGroupsWithTargets
}

var _ ddb.QueryOutputUnmarshalerWithPagination = &ListRequestWithGroupsWithTargets{}

func (g *ListRequestWithGroupsWithTargets) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := &dynamodb.QueryInput{
		KeyConditionExpression: aws.String("PK = :pk1"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk1": &types.AttributeValueMemberS{Value: keys.AccessRequest.PK1},
		},
	}

	return qi, nil
}

func (g *ListRequestWithGroupsWithTargets) UnmarshalQueryOutputWithPagination(out *dynamodb.QueryOutput) (map[string]types.AttributeValue, error) {
	result, pagination, err := UnmarshalRequests(out.Items)
	if err != nil {
		return nil, err
	}
	g.Result = result
	return pagination, nil
}
