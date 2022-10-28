package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/granted-approvals/pkg/access"
	"github.com/common-fate/granted-approvals/pkg/storage/keys"
)

type ListBookmarksForUser struct {
	UserID string
	Result []access.Bookmark `ddb:"result"`
}

func (l *ListBookmarksForUser) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := dynamodb.QueryInput{
		KeyConditionExpression: aws.String("PK = :pk1 and begins_with(SK, :sk1)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk1": &types.AttributeValueMemberS{Value: keys.Bookmark.PK1},
			":sk1": &types.AttributeValueMemberS{Value: keys.Bookmark.SK1User(l.UserID)},
		},
	}
	return &qi, nil
}
