package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/common-fate/pkg/target"
)

type ListValidTargetRoutesForGroupByPriority struct {
	Group  string
	Result []target.Route `ddb:"result"`
}

func (l *ListValidTargetRoutesForGroupByPriority) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := dynamodb.QueryInput{
		IndexName:              &keys.IndexNames.GSI1,
		ScanIndexForward:       aws.Bool(false),
		KeyConditionExpression: aws.String("GSI1PK = :pk and begins_with(GSI1SK, :sk)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: keys.TargetRoute.GSI1PK(l.Group)},
			":sk": &types.AttributeValueMemberS{Value: keys.TargetRoute.GSI1SKValid(true)},
		},
	}
	return &qi, nil
}
