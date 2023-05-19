package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/storage/keys"
)

type ListAccessTemplate struct {
	Result []access.AccessTemplate `ddb:"result"`
}

func (g *ListAccessTemplate) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := &dynamodb.QueryInput{
		KeyConditionExpression: aws.String("PK = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: keys.AccessTemplate.PK1},
		},
	}
	return qi, nil
}

// func (g *ListAccessTemplate) UnmarshalQueryOutput(out *dynamodb.QueryOutput) (*ddb.UnmarshalResult, error) {
// 	if len(out.Items) != 1 {
// 		return nil, ddb.ErrNoItems
// 	}

// 	return &ddb.UnmarshalResult{}, attributevalue.UnmarshalMap(out.Items[0], &g.Result)
// }
