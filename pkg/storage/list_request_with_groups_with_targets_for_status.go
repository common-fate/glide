package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	ddbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
)

type ListRequestWithGroupsWithTargetsForStatus struct {
	Status types.RequestStatus
	Result []access.RequestWithGroupsWithTargets
}

var _ ddb.QueryOutputUnmarshalerWithPagination = &ListRequestWithGroupsWithTargetsForStatus{}

func (g *ListRequestWithGroupsWithTargetsForStatus) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := &dynamodb.QueryInput{
		ScanIndexForward:       aws.Bool(false),
		IndexName:              &keys.IndexNames.GSI2,
		KeyConditionExpression: aws.String("GSI2PK = :pk1"),
		ExpressionAttributeValues: map[string]ddbTypes.AttributeValue{
			":pk1": &ddbTypes.AttributeValueMemberS{Value: keys.AccessRequest.GSI2PK(g.Status)},
		},
	}

	return qi, nil
}

func (g *ListRequestWithGroupsWithTargetsForStatus) UnmarshalQueryOutputWithPagination(out *dynamodb.QueryOutput) (map[string]ddbTypes.AttributeValue, error) {
	result, pagination, err := UnmarshalRequestsBottomToTop(out.Items)
	if err != nil {
		return nil, err
	}
	g.Result = result
	return pagination, nil
}
