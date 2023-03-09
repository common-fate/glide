package storage

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/common-fate/common-fate/pkg/storage/keys"
)

type ListUsersByFirstName struct {
	Result    []identity.User `ddb:"result"`
	FirstName string
}

func (l *ListUsersByFirstName) BuildQuery() (*dynamodb.QueryInput, error) {

	qi := dynamodb.QueryInput{
		IndexName:              aws.String(keys.IndexNames.GSI3),
		KeyConditionExpression: aws.String("GSI3PK = :pk and begins_with(GSI3SK, :sk1)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":  &types.AttributeValueMemberS{Value: keys.Users.GSI3PK},
			":sk1": &types.AttributeValueMemberS{Value: "ACTIVE#" + strings.ToLower(l.FirstName)},
		},
	}

	return &qi, nil
}
