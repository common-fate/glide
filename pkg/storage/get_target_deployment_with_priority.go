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

	Result targetgroup.Deployment `ddb:"result"`
}

func (g *GetTargetGroupDeploymentWithPriority) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := dynamodb.QueryInput{
		ScanIndexForward: aws.Bool(false),

		IndexName:              aws.String(keys.IndexNames.GSI1),
		Limit:                  aws.Int32(1),
		KeyConditionExpression: aws.String("GSI1PK = :pk and begins_with(GSI1SK, :sk)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: keys.TargetGroupDeployment.GSIPK1ValidHealthy(g.TargetGroupId)},
			":sk": &types.AttributeValueMemberS{Value: keys.TargetGroupDeployment.GSISK1ValidHealthy},
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
