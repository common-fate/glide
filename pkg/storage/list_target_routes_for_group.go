package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/common-fate/pkg/target"
)

type ListTargetRoutesForGroup struct {
	Group  string
	Result []target.Route `ddb:"result"`
}

func (l *ListTargetRoutesForGroup) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := dynamodb.QueryInput{
		KeyConditionExpression: aws.String("PK = :pk and begins_with(SK, :sk)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: keys.TargetRoute.PK1},
			":sk": &types.AttributeValueMemberS{Value: keys.TargetRoute.SK1Group(l.Group)},
		},
	}
	return &qi, nil
}
