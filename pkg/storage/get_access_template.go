package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/ddb"
)

type GetAccessTemplate struct {
	ID     string
	UserId string
	Result *access.AccessTemplate
}

func (g *GetAccessTemplate) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := &dynamodb.QueryInput{
		Limit:                  aws.Int32(1),
		KeyConditionExpression: aws.String("PK = :pk AND SK = :sk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: keys.AccessTemplate.PK1},
			":sk": &types.AttributeValueMemberS{Value: keys.AccessTemplate.SK1(g.ID)},
		},
	}
	return qi, nil
}
func (g *GetAccessTemplate) UnmarshalQueryOutput(out *dynamodb.QueryOutput) (*ddb.UnmarshalResult, error) {
	if len(out.Items) != 1 {
		return nil, ddb.ErrNoItems
	}

	return &ddb.UnmarshalResult{}, attributevalue.UnmarshalMap(out.Items[0], &g.Result)
}
