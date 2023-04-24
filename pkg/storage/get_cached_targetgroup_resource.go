package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/common-fate/pkg/cache"
	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/ddb"
)

type GetCachedTargetGroupResource struct {
	TargetGroupID string
	ResourceType  string
	ResourceID    string
	Result        *cache.TargetGroupResource
}

func (q *GetCachedTargetGroupResource) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := dynamodb.QueryInput{
		KeyConditionExpression: aws.String("PK = :pk1 and begins_with(SK, :sk1)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk1": &types.AttributeValueMemberS{Value: keys.TargetGroupResource.PK1},
			":sk1": &types.AttributeValueMemberS{Value: keys.TargetGroupResource.SK1(q.TargetGroupID, q.ResourceType, q.ResourceID)},
		},
	}
	return &qi, nil
}

func (q *GetCachedTargetGroupResource) UnmarshalQueryOutput(out *dynamodb.QueryOutput) (*ddb.UnmarshalResult, error) {
	if len(out.Items) == 0 {
		return nil, ddb.ErrNoItems
	}

	return &ddb.UnmarshalResult{}, attributevalue.UnmarshalMap(out.Items[0], &q.Result)

}
