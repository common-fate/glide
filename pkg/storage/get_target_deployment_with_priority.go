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

type GetTargetGroupDeploymentWithPriority struct {
	TargetGroupId string
	Valid         string
	Health        string
	Result        targetgroup.Deployment `ddb:"result"`
}

func (g *GetTargetGroupDeploymentWithPriority) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := dynamodb.QueryInput{
		IndexName:              aws.String(keys.IndexNames.GSI1),
		Limit:                  aws.Int32(1),
		KeyConditionExpression: aws.String("GSIPK = :pk and begins_with(GSISK, :sk1)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: keys.TargetGroupDeployment.GSIPK1(g.TargetGroupId)},
			":sk": &types.AttributeValueMemberS{Value: keys.TargetGroupDeployment.GSISK1(g.Valid, g.Health, "")},
			//where sk = true#true#
			//Will this ^ return the highest priority given the above query?
			//Will be saved to the database like this true#true#999
		},
	}
	return &qi, nil
}

func (g *GetTargetGroupDeploymentWithPriority) UnmarshalQueryOutput(out *dynamodb.QueryOutput) error {
	if len(out.Items) != 1 {
		return ddb.ErrNoItems
	}

	return attributevalue.UnmarshalMap(out.Items[0], &g.Result)
}
