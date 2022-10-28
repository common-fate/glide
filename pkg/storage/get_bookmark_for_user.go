package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/ddb"
	"github.com/common-fate/granted-approvals/pkg/access"
	"github.com/common-fate/granted-approvals/pkg/storage/keys"
)

type GetBookmarkForUser struct {
	ID     string
	UserID string
	Result *access.Bookmark
}

func (b *GetBookmarkForUser) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := &dynamodb.QueryInput{
		Limit:                  aws.Int32(1),
		KeyConditionExpression: aws.String("PK = :pk1 and SK = :sk1"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk1": &types.AttributeValueMemberS{Value: keys.Bookmark.PK1},
			":sk1": &types.AttributeValueMemberS{Value: keys.Bookmark.SK1(b.UserID, b.ID)},
		},
	}

	return qi, nil
}

func (g *GetBookmarkForUser) UnmarshalQueryOutput(out *dynamodb.QueryOutput) error {
	if len(out.Items) != 1 {
		return ddb.ErrNoItems
	}

	return attributevalue.UnmarshalMap(out.Items[0], &g.Result)
}
