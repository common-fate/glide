package ddb

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"

	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/ddb"
)

// Needed to be able to delete without duplicates
type UserBaseKeyOnly struct {
	ID string
}

func (u *UserBaseKeyOnly) DDBKeys() (ddb.Keys, error) {
	keys := ddb.Keys{
		PK: keys.Users.PK1,
		SK: keys.Users.SK1(u.ID),
	}

	return keys, nil
}

// Needed to filter by duplicates. Not needed in normal operations
type ListUsersForEmail struct {
	Result []identity.User `ddb:"result"`
	Email  string
}

func (l *ListUsersForEmail) BuildQuery() (*dynamodb.QueryInput, error) {

	qi := dynamodb.QueryInput{
		IndexName:              aws.String(keys.IndexNames.GSI2),
		KeyConditionExpression: aws.String("GSI2PK = :pk2 and GSI2SK = :sk2"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk2": &types.AttributeValueMemberS{Value: keys.Users.GSI2PK},
			":sk2": &types.AttributeValueMemberS{Value: keys.Users.GSI2SK(string(l.Email))},
		},
	}
	return &qi, nil
}
