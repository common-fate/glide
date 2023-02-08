package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/common-fate/pkg/targetgroup"
	"github.com/common-fate/ddb"
)

type GetTargetGroupDeployment struct {
	ID     string
	Result targetgroup.Deployment `ddb:"result"`
}

func (g *GetTargetGroupDeployment) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := dynamodb.QueryInput{
		Limit:                  aws.Int32(1),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :sk)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: keys.TargetGroupDeployment.PK1},
			":sk": &types.AttributeValueMemberS{Value: keys.TargetGroupDeployment.SK1(g.ID)},
		},
	}
	return &qi, nil
}

func (g *GetTargetGroupDeployment) UnmarshalQueryOutput(out *dynamodb.QueryOutput) error {
	if len(out.Items) != 1 {
		return ddb.ErrNoItems
	}

	return attributevalue.UnmarshalMap(out.Items[0], &g.Result)
}
