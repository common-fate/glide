package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/storage/keys"
)

type ListRequestsForUser struct {
	UserId string
	Result []access.Request `ddb:"result"`
}

func (g *ListRequestsForUser) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := dynamodb.QueryInput{
		// newest to oldest
		ScanIndexForward:       aws.Bool(false),
		IndexName:              aws.String(keys.IndexNames.GSI1),
		KeyConditionExpression: aws.String("GSI1PK = :pk "),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: keys.AccessRequest.GSI1PK(g.UserId)},
		},
	}
	return &qi, nil
}
