package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/ddb"
)

type ListCachedProviderV2Resource struct {
	ProviderID   string
	ResourceType string
	ResourceId   string
	Resource     interface{}
}

// TODO: Maybe we need to create a new SK for resourceType
func (q *ListCachedProviderV2Resource) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := dynamodb.QueryInput{
		KeyConditionExpression: aws.String("PK = :pk1 and begins_with(SK, :sk1)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk1": &types.AttributeValueMemberS{Value: keys.ProviderOption.PK1},
		},
	}
	return &qi, nil
}

func (q *ListCachedProviderV2Resource) UnmarshalQueryOutput(out *dynamodb.QueryOutput) error {
	if len(out.Items) == 0 {
		return ddb.ErrNoItems
	}

	return attributevalue.UnmarshalListOfMaps(out.Items, &q.Resource)
}
