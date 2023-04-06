package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/common-fate/pkg/requestsv2.go"
	"github.com/common-fate/common-fate/pkg/storage/keys"
)

type GetRequestGroup struct {
	ID     string
	Result *requestsv2.RequestGroup `ddb:"result"`
}

func (g *GetRequestGroup) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := dynamodb.QueryInput{
		Limit:                  aws.Int32(1),
		KeyConditionExpression: aws.String("PK = :pk AND PK = :sk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: keys.RequestGroups.PK1},
			":sk": &types.AttributeValueMemberS{Value: keys.RequestGroups.SK1(g.ID)},
		},
	}
	return &qi, nil
}
