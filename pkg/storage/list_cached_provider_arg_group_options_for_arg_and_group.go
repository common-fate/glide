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

type ListCachedProviderArgGroupOptionsForArgAndGroup struct {
	ProviderID string
	ArgID      string
	GroupID    string
	Result     []cache.ProviderArgGroupOption
}

func (q *ListCachedProviderArgGroupOptionsForArgAndGroup) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := dynamodb.QueryInput{
		KeyConditionExpression: aws.String("PK = :pk1 and begins_with(SK, :sk1)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk1": &types.AttributeValueMemberS{Value: keys.ProviderArgGroupOption.PK1},
			":sk1": &types.AttributeValueMemberS{Value: keys.ProviderArgGroupOption.SK1ProviderArg(q.ProviderID, q.ArgID)},
		},
	}
	return &qi, nil
}

func (q *ListCachedProviderArgGroupOptionsForArgAndGroup) UnmarshalQueryOutput(out *dynamodb.QueryOutput) error {
	if len(out.Items) == 0 {
		return ddb.ErrNoItems
	}

	return attributevalue.UnmarshalListOfMaps(out.Items, &q.Result)

}
