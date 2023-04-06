package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/common-fate/pkg/requestsv2.go"
	"github.com/common-fate/common-fate/pkg/storage/keys"
)

type GetRequestV2 struct {
	ID     string
	UserId string
	Result *requestsv2.Requestv2
}

func (g *GetRequestV2) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := &dynamodb.QueryInput{
		Limit:                  aws.Int32(1),
		KeyConditionExpression: aws.String("PK = :pk AND SK = :sk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: keys.RequestV2.PK1},
			":sk": &types.AttributeValueMemberS{Value: keys.RequestV2.SK1(g.UserId, g.ID)},
		},
	}
	return qi, nil
}
