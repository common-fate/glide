package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/storage/keys"
)

type ListRequestsForStatus struct {
	Status access.Status
	Result []access.Request `ddb:"result"`
}

func (l *ListRequestsForStatus) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := dynamodb.QueryInput{
		// newest to oldest
		ScanIndexForward:       aws.Bool(false),
		IndexName:              aws.String(keys.IndexNames.GSI2),
		KeyConditionExpression: aws.String("GSI2PK = :pk1"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk1": &types.AttributeValueMemberS{Value: keys.AccessRequest.GSI2PK(string(l.Status))},
		},
	}
	return &qi, nil
}
