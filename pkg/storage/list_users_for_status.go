package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/granted-approvals/pkg/identity"
	"github.com/common-fate/granted-approvals/pkg/storage/keys"
	gendTypes "github.com/common-fate/granted-approvals/pkg/types"
)

type ListUsersForStatus struct {
	Result []identity.User `ddb:"result"`
	Status gendTypes.IdpStatus
}

func (l *ListUsersForStatus) BuildQuery() (*dynamodb.QueryInput, error) {

	qi := dynamodb.QueryInput{
		IndexName:              aws.String(keys.IndexNames.GSI1),
		KeyConditionExpression: aws.String("GSI1PK = :pk1 and begins_with(GSI1SK, :sk1)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk1": &types.AttributeValueMemberS{Value: keys.Users.GSI1PK},
			":sk1": &types.AttributeValueMemberS{Value: keys.Users.GSI1SKStatus(string(l.Status))},
		},
	}
	return &qi, nil
}
