package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/common-fate/pkg/requests"
	"github.com/common-fate/common-fate/pkg/storage/keys"
)

type ListRequestV2 struct {
	ID     string
	UserId string
	Result []requests.Requestv2
}

func (g *ListRequestV2) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := &dynamodb.QueryInput{
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(PK, :sk)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: keys.RequestV2.PK1},
			":sk": &types.AttributeValueMemberS{Value: keys.RequestV2.SKAllRequests(g.UserId)},
		},
	}
	return qi, nil
}
