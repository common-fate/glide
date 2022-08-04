package fixtures

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers/azure/ad"
	"github.com/common-fate/granted-approvals/pkg/gconfig"
	"github.com/segmentio/ksuid"
)

type Fixtures struct {
	User    string
	GroupID string
}

type Generator struct {
	client ad.Provider
}

func (a *Generator) Config() gconfig.Config {
	return a.client.Config()
}

// Init the Azure provider.
func (a *Generator) Init(ctx context.Context) error {
	return a.client.Init(ctx)
}

// Generate fixtures by calling the azure API.
func (g *Generator) Generate(ctx context.Context) ([]byte, error) {
	email := fmt.Sprintf("test_%s@exponentlabsio.onmicrosoft.com", ksuid.New().String())

	createUserRequest := ad.CreateADUser{
		AccountEnabled:    true,
		DisplayName:       "test",
		MailNickname:      "test",
		UserPrincipalName: email,
		PasswordProfile:   ad.PasswordProfile{ForceChangePasswordNextSignIn: true, Password: ksuid.New().String()},
	}

	err := g.client.CreateUser(ctx, createUserRequest)
	if err != nil {
		return nil, err
	}

	group := fmt.Sprintf("group_%s", ksuid.New().String())

	creategroup := ad.CreateADGroup{
		Description:     "test group",
		DisplayName:     group,
		MailNickname:    fmt.Sprintf("test_%s", ksuid.New().String()),
		GroupTypes:      []string{"Unified"},
		SecurityEnabled: false,
		MailEnabled:     true,
	}

	azureGroup, err := g.client.CreateGroup(ctx, creategroup)
	if err != nil {
		return nil, err
	}

	f := Fixtures{
		GroupID: azureGroup.ID,
		User:    email,
	}

	return json.Marshal(f)
}

func (g *Generator) Destroy(ctx context.Context, data []byte) error {
	var f Fixtures
	err := json.Unmarshal(data, &f)
	if err != nil {
		return err
	}

	err = g.client.DeleteGroup(ctx, f.GroupID)
	if err != nil {
		return err
	}

	err = g.client.DeleteUser(ctx, f.User)
	return err
}
