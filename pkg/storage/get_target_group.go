package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/common-fate/pkg/targetgroup"
)

type GetTargetGroup struct {
	ID     string
	Result targetgroup.TargetGroup `ddb:"result"`
}

func (g *GetTargetGroup) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := dynamodb.QueryInput{
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :sk)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: keys.TargetGroup.PK1},
			":sk": &types.AttributeValueMemberS{Value: keys.TargetGroup.SK1(g.ID)},
		},
	}
	return &qi, nil
}
