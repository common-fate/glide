package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/storage/keys"
)

type ListRequestEvents struct {
	RequestID string
	Result    []access.RequestEvent `ddb:"result"`
}

func (l *ListRequestEvents) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := dynamodb.QueryInput{
		KeyConditionExpression: aws.String("PK = :pk1 AND begins_with(SK, :sk1)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk1": &types.AttributeValueMemberS{Value: keys.AccessRequestEvent.PK1},
			":sk1": &types.AttributeValueMemberS{Value: keys.AccessRequestEvent.SK1Request(l.RequestID)},
		},
	}
	return &qi, nil
}
