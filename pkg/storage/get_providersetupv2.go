package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	cf_types "github.com/common-fate/common-fate/pkg/types"

	"github.com/common-fate/common-fate/pkg/providersetupv2"
	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/ddb"
)

type GetProviderSetupV2 struct {
	ID     string
	Result *providersetupv2.Setup
}

func (g *GetProviderSetupV2) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := &dynamodb.QueryInput{
		Limit:                  aws.Int32(1),
		KeyConditionExpression: aws.String("PK = :pk AND SK = :sk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: keys.ProviderSetupV2.PK1},
			":sk": &types.AttributeValueMemberS{Value: keys.ProviderSetupV2.SK1(g.ID)},
		},
	}
	return qi, nil
}

func (g *GetProviderSetupV2) UnmarshalQueryOutput(out *dynamodb.QueryOutput) error {
	if len(out.Items) != 1 {
		return ddb.ErrNoItems
	}

	return attributevalue.UnmarshalMap(out.Items[0], &g.Result)
}

func (g *GetProviderSetupV2) ToAPI() cf_types.ProviderSetupV2 {
	return cf_types.ProviderSetupV2{
		Id:           g.Result.ID,
		Name:         g.Result.ProviderName,
		Team:         g.Result.ProviderTeam,
		Status:       g.Result.Status,
		Version:      g.Result.ProviderVersion,
		ConfigValues: g.Result.ConfigValues,
		// Steps:        g.Result.Steps,
	}

	//todo manually loop through steps and configvalidation

}
