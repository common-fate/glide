package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/common-fate/pkg/targetgroup"
)

type ListDeploymentsForTargetGroup struct {
	TargetGroupId string
	Result        []targetgroup.Deployment `ddb:"result"`
}

func (g *ListDeploymentsForTargetGroup) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := dynamodb.QueryInput{
		IndexName:              aws.String(keys.IndexNames.GSI1),
		KeyConditionExpression: aws.String("GSIPK = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: keys.TargetGroupDeployment.GSIPK1(g.TargetGroupId)},
		},
	}
	return &qi, nil
}
