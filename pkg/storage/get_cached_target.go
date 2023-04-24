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

type GetCachedTarget struct {
	ID     string
	Result *cache.Target `ddb:"result"`
}

func (l *GetCachedTarget) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := dynamodb.QueryInput{
		KeyConditionExpression: aws.String("PK = :pk and SK = :sk"),
		Limit:                  aws.Int32(1),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: keys.EntitlementTarget.PK1},
			":sk": &types.AttributeValueMemberS{Value: keys.EntitlementTarget.SK1(l.ID)},
		},
	}
	return &qi, nil
}
func (q *GetCachedTarget) UnmarshalQueryOutput(out *dynamodb.QueryOutput) (*ddb.UnmarshalResult, error) {
	if len(out.Items) != 1 {
		return nil, ddb.ErrNoItems
	}
	return &ddb.UnmarshalResult{}, attributevalue.UnmarshalMap(out.Items[0], &q.Result)
}
