package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/granted-approvals/pkg/identity"
	"github.com/common-fate/granted-approvals/pkg/storage/keys"
	gendTypes "github.com/common-fate/granted-approvals/pkg/types"
)

type ListGroupsForSourceAndStatus struct {
	Result []identity.Group `ddb:"result"`
	Source string
	Status gendTypes.IdpStatus
}

func (l *ListGroupsForSourceAndStatus) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := dynamodb.QueryInput{
		IndexName: aws.String(keys.IndexNames.GSI2),

		KeyConditionExpression: aws.String("GSI2PK = :pk1 and GSI2SK = :sk1"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk1": &types.AttributeValueMemberS{Value: keys.Groups.GSI2PK},
			":sk1": &types.AttributeValueMemberS{Value: keys.Groups.GSI2SKSourceStatus(l.Source, string(l.Status))},
		},
	}
	return &qi, nil
}
