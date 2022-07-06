package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/granted-approvals/pkg/access"
	"github.com/common-fate/granted-approvals/pkg/storage/keys"
)

type ListRequestsForReviewerAndStatus struct {
	ReviewerID string
	Status     access.Status
	Result     []access.Request `ddb:"result"`
}

func (l *ListRequestsForReviewerAndStatus) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := dynamodb.QueryInput{
		// newest to oldest
		ScanIndexForward:       aws.Bool(false),
		IndexName:              aws.String(keys.IndexNames.GSI2),
		KeyConditionExpression: aws.String("GSI2PK = :pk and begins_with(GSI2SK, :sk)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: keys.RequestReviewer.GSI2PK(l.ReviewerID)},
			":sk": &types.AttributeValueMemberS{Value: keys.RequestReviewer.GSI2SKStatus(string(l.Status))},
		},
	}
	return &qi, nil
}
func (g *ListRequestsForReviewerAndStatus) UnmarshalQueryOutput(out *dynamodb.QueryOutput) error {
	var r []access.Reviewer
	err := attributevalue.UnmarshalListOfMaps(out.Items, &r)
	if err != nil {
		return err
	}
	for _, r2 := range r {
		g.Result = append(g.Result, r2.Request)
	}
	return nil
}
