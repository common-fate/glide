package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/ddb"
	"github.com/common-fate/granted-approvals/pkg/identity"
	"github.com/common-fate/granted-approvals/pkg/storage/keys"
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
func (g *GetUserByEmail) UnmarshalQueryOutput(out *dynamodb.QueryOutput) error {
	if len(out.Items) != 1 {
		return ddb.ErrNoItems
	}

	return attributevalue.UnmarshalMap(out.Items[0], &g.Result)
}
