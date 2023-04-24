package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/ddb"
)

type GetUserByEmail struct {
	Email  string
	Result *identity.User
}

func (u *GetUserByEmail) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := &dynamodb.QueryInput{
		Limit:                  aws.Int32(1),
		IndexName:              &keys.IndexNames.GSI2,
		KeyConditionExpression: aws.String("GSI2PK = :pk1 and GSI2SK = :sk1"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk1": &types.AttributeValueMemberS{Value: keys.Users.GSI2PK},
			":sk1": &types.AttributeValueMemberS{Value: keys.Users.GSI2SK(u.Email)},
		},
	}

	return qi, nil
}
func (g *GetUserByEmail) UnmarshalQueryOutput(out *dynamodb.QueryOutput) (*ddb.UnmarshalResult, error) {
	if len(out.Items) != 1 {
		return nil, ddb.ErrNoItems
	}

	return &ddb.UnmarshalResult{}, attributevalue.UnmarshalMap(out.Items[0], &g.Result)
}
