package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/ddb"
)

type GetFavoriteForUser struct {
	ID     string
	UserID string
	Result *access.Favorite
}

func (b *GetFavoriteForUser) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := &dynamodb.QueryInput{
		Limit:                  aws.Int32(1),
		KeyConditionExpression: aws.String("PK = :pk1 and SK = :sk1"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk1": &types.AttributeValueMemberS{Value: keys.Favorite.PK1},
			":sk1": &types.AttributeValueMemberS{Value: keys.Favorite.SK1(b.UserID, b.ID)},
		},
	}

	return qi, nil
}

func (g *GetFavoriteForUser) UnmarshalQueryOutput(out *dynamodb.QueryOutput) error {
	if len(out.Items) != 1 {
		return ddb.ErrNoItems
	}

	return attributevalue.UnmarshalMap(out.Items[0], &g.Result)
}
