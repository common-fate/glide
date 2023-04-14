package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/common-fate/pkg/requests"
	"github.com/common-fate/common-fate/pkg/storage/keys"
)

type GetGrantV2 struct {
	GroupID string
	GrantId string
	Result  *requests.Grantv2
}

func (g *ListGrantsV2) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := &dynamodb.QueryInput{
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(PK, :sk)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: keys.Grant.PK1},
			":sk": &types.AttributeValueMemberS{Value: keys.Grant.SKAllGrants(g.GroupID)},
		},
	}
	return qi, nil
}
