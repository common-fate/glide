package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/ddb"
	"github.com/common-fate/granted-approvals/pkg/cache"
	"github.com/common-fate/granted-approvals/pkg/storage/keys"
)

type ListCachedProviderArgGroupOptionValueForArg struct {
	ProviderID string
	ArgID      string
	GroupId    string
	GroupValue string
	Result     []cache.ProviderArgGroupOption
}

func (q *ListCachedProviderArgGroupOptionValueForArg) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := dynamodb.QueryInput{
		KeyConditionExpression: aws.String("PK = :pk1 and begins_with(SK, :sk1)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk1": &types.AttributeValueMemberS{Value: keys.ProviderArgGroupOption.PK1},
			":sk1": &types.AttributeValueMemberS{Value: keys.ProviderArgGroupOption.SK1(q.ProviderID, q.ArgID, q.GroupId, q.GroupValue)},
		},
	}
	return &qi, nil
}

func (q *ListCachedProviderArgGroupOptionValueForArg) UnmarshalQueryOutput(out *dynamodb.QueryOutput) error {
	if len(out.Items) == 0 {
		return ddb.ErrNoItems
	}

	return attributevalue.UnmarshalListOfMaps(out.Items, &q.Result)

}
