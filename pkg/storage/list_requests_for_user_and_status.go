package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/granted-approvals/pkg/access"
	"github.com/common-fate/granted-approvals/pkg/storage/keys"
)

type ListRequestsForUserAndStatus struct {
	Status access.Status
	UserId string
	Result []access.Request `ddb:"result"`
}

func (g *ListRequestsForUserAndStatus) BuildQuery() (*dynamodb.QueryInput, error) {

	qi := dynamodb.QueryInput{
		// newest to oldest
		ScanIndexForward:       aws.Bool(false),
		IndexName:              aws.String(keys.IndexNames.GSI2),
		KeyConditionExpression: aws.String("GSI2PK = :pk and begins_with(GSI2SK, :sk)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: keys.AccessRequest.GSI2PK(string(g.Status))},
			":sk": &types.AttributeValueMemberS{Value: keys.AccessRequest.GSI2SKUser(g.UserId)},
		},
	}

	return &qi, nil
}
