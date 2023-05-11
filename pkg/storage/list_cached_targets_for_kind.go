package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/common-fate/pkg/cache"
	"github.com/common-fate/common-fate/pkg/storage/keys"
)

type ListCachedTargetsForKind struct {
	Publisher string
	Name      string
	Kind      string
	Result    []cache.Target `ddb:"result"`
}

func (l *ListCachedTargetsForKind) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := dynamodb.QueryInput{
		KeyConditionExpression: aws.String("PK = :pk and begins_with(SK, :sk)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: keys.EntitlementTarget.PK1},
			":sk": &types.AttributeValueMemberS{Value: keys.EntitlementTarget.SK1PublisherNameKind(l.Publisher, l.Name, l.Kind)},
		},
	}
	return &qi, nil
}
