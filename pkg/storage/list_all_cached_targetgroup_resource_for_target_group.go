package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/common-fate/pkg/cache"
	"github.com/common-fate/common-fate/pkg/storage/keys"
)

type ListCachedTargetGroupResourceForTargetGroupAndResourceType struct {
	TargetGroupID string
	ResourceType  string
	Result        []cache.TargetGroupResource `ddb:"result"`
}

func (q *ListCachedTargetGroupResourceForTargetGroupAndResourceType) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := dynamodb.QueryInput{
		KeyConditionExpression: aws.String("PK = :pk1 and begins_with(SK, :sk1)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk1": &types.AttributeValueMemberS{Value: keys.TargetGroupResource.PK1},
			":sk1": &types.AttributeValueMemberS{Value: keys.TargetGroupResource.SK1TargetGroupResource(q.TargetGroupID, q.ResourceType)},
		},
	}
	return &qi, nil
}
