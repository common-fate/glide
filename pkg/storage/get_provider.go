package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/provider-registry/pkg/provider"
)

type GetProvider struct {
	ID     string
	Result *provider.Provider
}

func (g *GetProvider) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := &dynamodb.QueryInput{
		Limit:                  aws.Int32(1),
		KeyConditionExpression: aws.String("PK = :pk1 and SK = :sk1"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk1": &types.AttributeValueMemberS{Value: keys.Provider.PK1},
			":sk1": &types.AttributeValueMemberS{Value: keys.Provider.SK1(g.ID)},
		},
	}

	return qi, nil
}

// func (g *GetProvider) UnmarshalQueryOutput(out *dynamodb.QueryOutput) error {
// 	if len(out.Items) != 1 {
// 		return ddb.ErrNoItems
// 	}
// 	return UnmarshalMap(out.Items[0], &g.Result)
// }
