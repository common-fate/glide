package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/storage/keys"
)

type ListRequestsForReviewer struct {
	ReviewerID string
	Result     []access.Request `ddb:"result"`
}

func (l *ListRequestsForReviewer) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := dynamodb.QueryInput{
		// newest to oldest
		ScanIndexForward:       aws.Bool(false),
		IndexName:              aws.String(keys.IndexNames.GSI1),
		KeyConditionExpression: aws.String("GSI1PK = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: keys.RequestReviewer.GSI1PK(l.ReviewerID)},
		},
	}
	return &qi, nil
}
func (g *ListRequestsForReviewer) UnmarshalQueryOutput(out *dynamodb.QueryOutput) error {
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
