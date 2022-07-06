package fixtures

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssoadmin"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/genv"
	"github.com/segmentio/ksuid"
)

type Fixtures struct {
	PermissionSetARN string
	AccountID        string
	User             string
}

type Generator struct {
	client      *ssoadmin.Client
	InstanceARN string
	AccountID   string
	User        string
}

// Configure the fixture generator
func (g *Generator) Config() genv.Config {
	return genv.Config{
		genv.String("instanceArn", &g.InstanceARN, ""),
		genv.SecretString("fixturesAccountId", &g.AccountID, ""),
		genv.SecretString("fixturesUser", &g.User, ""),
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
		InstanceArn: &g.InstanceARN,
		Name:        &name,
		Description: aws.String("Granted Integration Testing"),
	})

	if err != nil {
		return nil, err
	}

	f := Fixtures{
		PermissionSetARN: *res.PermissionSet.PermissionSetArn,
		AccountID:        g.AccountID,
		User:             g.User,
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
		InstanceArn:      &g.InstanceARN,
		PermissionSetArn: &f.PermissionSetARN,
	})
	return err
}
