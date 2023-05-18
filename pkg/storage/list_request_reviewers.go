package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/storage/keys"
)

type ListRequestReviewers struct {
	RequestId string
	Result    []access.Reviewer `ddb:"result"`
}

func (g *ListRequestReviewers) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := dynamodb.QueryInput{
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :sk)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: keys.RequestReviewer.PK1},
			":sk": &types.AttributeValueMemberS{Value: g.RequestId},
		},
	}
	return &qi, nil
}
