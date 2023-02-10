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

type GetTargetGroupDeploymentWithGroupId struct {
	TargetGroupId string
	Result        targetgroup.Deployment `ddb:"result"`
}

func (g *GetTargetGroupDeploymentWithGroupId) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := dynamodb.QueryInput{
		IndexName:              aws.String(keys.IndexNames.GSI1),
		Limit:                  aws.Int32(1),
		KeyConditionExpression: aws.String("GSIPK = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: keys.TargetGroupDeployment.PK1},
		},
	}
	return &qi, nil
}

func (g *GetTargetGroupDeploymentWithGroupId) UnmarshalQueryOutput(out *dynamodb.QueryOutput) error {
	if len(out.Items) != 1 {
		return ddb.ErrNoItems
	}

	return attributevalue.UnmarshalMap(out.Items[0], &g.Result)
}
