package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/common-fate/pkg/cache"
	"github.com/common-fate/common-fate/pkg/storage/keys"
)

type ListCachedTargetGroupResources struct {
	Result []cache.TargetGroupResource `ddb:"result"`
}

func (q *ListCachedTargetGroupResources) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := dynamodb.QueryInput{
		KeyConditionExpression: aws.String("PK = :pk1"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk1": &types.AttributeValueMemberS{Value: keys.TargetGroupResource.PK1},
		},
	}
	return &qi, nil
}
