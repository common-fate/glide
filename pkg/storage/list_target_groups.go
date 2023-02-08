package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/common-fate/pkg/targetgroup"
)

type ListTargetGroups struct {
	Result []targetgroup.TargetGroup `ddb:"result"`
}

func (l *ListTargetGroups) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := dynamodb.QueryInput{
		KeyConditionExpression: aws.String("PK = :pk "),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: keys.TargetGroup.PK1},
		},
	}
	return &qi, nil
}
