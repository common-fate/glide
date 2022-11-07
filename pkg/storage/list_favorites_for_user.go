package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/granted-approvals/pkg/access"
	"github.com/common-fate/granted-approvals/pkg/storage/keys"
)

type ListFavoritesForUser struct {
	UserID string
	Result []access.Favorite `ddb:"result"`
}

func (l *ListFavoritesForUser) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := dynamodb.QueryInput{
		KeyConditionExpression: aws.String("PK = :pk1 and begins_with(SK, :sk1)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk1": &types.AttributeValueMemberS{Value: keys.Favorite.PK1},
			":sk1": &types.AttributeValueMemberS{Value: keys.Favorite.SK1User(l.UserID)},
		},
	}
	return &qi, nil
}
