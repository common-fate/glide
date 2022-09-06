package deploy

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/ddb"
	"github.com/common-fate/granted-approvals/pkg/storage/keys"
)

// DDBManagedDeploymentConfig reads config values from DynamoDB
type DDBManagedDeploymentConfig struct {
	DB ddb.Storage
}

func (d *DDBManagedDeploymentConfig) ReadProviders(ctx context.Context) (ProviderMap, error) {
	q := getProviderConfig{}
	_, err := d.DB.Query(ctx, &q)
	if err != nil {
		return ProviderMap{}, err
	}
	if q.Result != nil {
		return q.Result.ProviderConfig, nil
	}
	return ProviderMap{}, nil
}

func (d *DDBManagedDeploymentConfig) ReadNotifications(ctx context.Context) (FeatureMap, error) {
	q := getNotificationConfig{}
	_, err := d.DB.Query(ctx, &q)
	if err != nil {
		return FeatureMap{}, err
	}
	if q.Result != nil {
		return q.Result.Notifications, nil
	}
	return FeatureMap{}, nil
}

func (d *DDBManagedDeploymentConfig) WriteProviders(ctx context.Context, pm ProviderMap) error {
	item := ddbProviderConfig{ProviderConfig: pm}
	return d.DB.Put(ctx, &item)
}

func (d *DDBManagedDeploymentConfig) WriteNotifications(ctx context.Context, fm FeatureMap) error {
	item := ddbNotificationConfig{Notifications: fm}
	return d.DB.Put(ctx, &item)
}

type ddbProviderConfig struct {
	ProviderConfig ProviderMap `json:"providerConfig" dynamodbav:"providerConfig"`
}

func (r *ddbProviderConfig) DDBKeys() (ddb.Keys, error) {
	keys := ddb.Keys{
		PK: keys.ManagedProviderConfig.PK,
		SK: keys.ManagedProviderConfig.SK,
	}

	return keys, nil
}

type getProviderConfig struct {
	Result *ddbProviderConfig
}

func (g *getProviderConfig) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := dynamodb.QueryInput{
		Limit:                  aws.Int32(1),
		KeyConditionExpression: aws.String("PK = :pk and SK = :sk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: keys.ManagedProviderConfig.PK},
			":sk": &types.AttributeValueMemberS{Value: keys.ManagedProviderConfig.PK},
		},
	}
	return &qi, nil
}

func (g *getProviderConfig) UnmarshalQueryOutput(out *dynamodb.QueryOutput) error {
	if len(out.Items) != 1 {
		g.Result = &ddbProviderConfig{}
		return nil
	}

	return attributevalue.UnmarshalMap(out.Items[0], &g.Result)
}

type ddbNotificationConfig struct {
	Notifications FeatureMap `json:"notifications" dynamodbav:"notifications"`
}

func (r *ddbNotificationConfig) DDBKeys() (ddb.Keys, error) {
	keys := ddb.Keys{
		PK: keys.ManagedNotificationConfig.PK,
		SK: keys.ManagedNotificationConfig.SK,
	}

	return keys, nil
}

type getNotificationConfig struct {
	Result *ddbNotificationConfig
}

func (g *getNotificationConfig) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := dynamodb.QueryInput{
		Limit:                  aws.Int32(1),
		KeyConditionExpression: aws.String("PK = :pk and SK = :sk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: keys.ManagedNotificationConfig.PK},
			":sk": &types.AttributeValueMemberS{Value: keys.ManagedNotificationConfig.PK},
		},
	}
	return &qi, nil
}

func (g *getNotificationConfig) UnmarshalQueryOutput(out *dynamodb.QueryOutput) error {
	if len(out.Items) != 1 {
		g.Result = &ddbNotificationConfig{}
		return nil
	}

	return attributevalue.UnmarshalMap(out.Items[0], &g.Result)
}
