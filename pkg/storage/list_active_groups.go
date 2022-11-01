package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/granted-approvals/pkg/identity"
	"github.com/common-fate/granted-approvals/pkg/storage/keys"
)

type ListActiveGroups struct {
	Result []identity.Group `ddb:"result"`
}

func (l *ListActiveGroups) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := dynamodb.QueryInput{
		KeyConditionExpression: aws.String("PK = :pk1"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk1": &types.AttributeValueMemberS{Value: keys.Groups.PK1},
		},
		ExpressionAttributeNames: make(map[string]string),
	}

	//Filter expression to not pull internal groups.
	var expr string

	expr += "#group_status <> :key"
	key := "ARCHIVED"
	qi.ExpressionAttributeValues[":key"] = &types.AttributeValueMemberS{Value: key}
	qi.ExpressionAttributeNames["#group_status"] = "status"
	qi.FilterExpression = &expr
	return &qi, nil
}
