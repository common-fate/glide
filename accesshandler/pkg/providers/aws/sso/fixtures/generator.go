package fixtures

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssoadmin"
	"github.com/common-fate/granted-approvals/pkg/gconfig"
	"github.com/segmentio/ksuid"
)

type Fixtures struct {
	PermissionSetARN string
	AccountID        string
	User             string
}

type Generator struct {
	client      *ssoadmin.Client
	InstanceARN gconfig.StringValue
	AccountID   gconfig.SecretStringValue
	User        gconfig.SecretStringValue
}

// Configure the fixture generator
func (g *Generator) Config() gconfig.Config {
	return gconfig.Config{
		Fields: []*gconfig.Field{
			gconfig.StringField("instanceArn", &g.InstanceARN, ""),
			gconfig.SecretStringField("fixturesAccountId", &g.AccountID, "", gconfig.WithNoArgs("")),
			gconfig.SecretStringField("fixturesUser", &g.User, "", gconfig.WithNoArgs("")),
		},
	}
}

func (g *Generator) Init(ctx context.Context) error {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}

	client := ssoadmin.NewFromConfig(cfg)

	g.client = client
	return nil
}

// Generate fixtures by calling the AWS SSO API.
func (g *Generator) Generate(ctx context.Context) ([]byte, error) {
	name := fmt.Sprintf("test%s", ksuid.New().String())
	res, err := g.client.CreatePermissionSet(ctx, &ssoadmin.CreatePermissionSetInput{
		InstanceArn: aws.String(g.InstanceARN.Get()),
		Name:        &name,
		Description: aws.String("Granted Integration Testing"),
	})

	if err != nil {
		return nil, err
	}

	f := Fixtures{
		PermissionSetARN: *res.PermissionSet.PermissionSetArn,
		AccountID:        g.AccountID.Get(),
		User:             g.User.Get(),
	}

	return json.Marshal(f)
}

func (g *Generator) Destroy(ctx context.Context, data []byte) error {
	var f Fixtures
	err := json.Unmarshal(data, &f)
	if err != nil {
		return err
	}

	_, err = g.client.DeletePermissionSet(ctx, &ssoadmin.DeletePermissionSetInput{
		InstanceArn:      aws.String(g.InstanceARN.Get()),
		PermissionSetArn: &f.PermissionSetARN,
	})
	return err
}
