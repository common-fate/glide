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

type GetCachedProviderArgGroupOptionValueForArg struct {
	ProviderID string
	ArgID      string
	GroupId    string
	GroupValue string
	Result     *cache.ProviderArgGroupOption
}

func (q *GetCachedProviderArgGroupOptionValueForArg) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := dynamodb.QueryInput{
		KeyConditionExpression: aws.String("PK = :pk1 and begins_with(SK, :sk1)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk1": &types.AttributeValueMemberS{Value: keys.ProviderArgGroupOption.PK1},
			":sk1": &types.AttributeValueMemberS{Value: keys.ProviderArgGroupOption.SK1(q.ProviderID, q.ArgID, q.GroupId, q.GroupValue)},
		},
	}
	return &qi, nil
}

func (q *GetCachedProviderArgGroupOptionValueForArg) UnmarshalQueryOutput(out *dynamodb.QueryOutput) error {
	if len(out.Items) != 1 {
		return ddb.ErrNoItems
	}
	return attributevalue.UnmarshalMap(out.Items[0], &q.Result)
}
