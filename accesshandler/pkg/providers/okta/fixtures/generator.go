package fixtures

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/common-fate/granted-approvals/pkg/gconfig"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
	"github.com/segmentio/ksuid"
)

type Fixtures struct {
	User    string
	GroupID string
}

type Generator struct {
	client   *okta.Client
	orgURL   gconfig.StringValue
	apiToken gconfig.SecretStringValue
}

// Configure the fixture generator
func (g *Generator) Config() gconfig.Config {
	return gconfig.Config{
		Fields: []*gconfig.Field{
			gconfig.StringField("orgUrl", &g.orgURL, "Okta org URL"),
			gconfig.SecretStringField("apiToken", &g.apiToken, "Okta API token", gconfig.WithArgs("/granted/providers/%s/apiToken", 1)),
		},
	}
}

func (g *Generator) Init(ctx context.Context) error {
	_, client, err := okta.NewClient(ctx, okta.WithOrgUrl(g.orgURL.Get()), okta.WithToken(g.apiToken.Get()))
	if err != nil {
		return err
	}

	g.client = client
	return nil
}

// Generate fixtures by calling the Okta API.
func (g *Generator) Generate(ctx context.Context) ([]byte, error) {
	email := fmt.Sprintf("test_%s@noreply.local", ksuid.New().String())

	createUserRequest := okta.CreateUserRequest{Profile: &okta.UserProfile{
		"email":     email,
		"login":     email,
		"firstName": "Test",
		"lastName":  "User",
	}}

	qp := query.NewQueryParams(query.WithActivate(true))

	_, _, err := g.client.User.CreateUser(ctx, createUserRequest, qp)
	if err != nil {
		return nil, err
	}

	group := fmt.Sprintf("group_%s", ksuid.New().String())

	oktaGroup, _, err := g.client.Group.CreateGroup(ctx, okta.Group{
		Profile: &okta.GroupProfile{
			Name: group,
		},
	})
	if err != nil {
		return nil, err
	}

	f := Fixtures{
		GroupID: oktaGroup.Id,
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

	_, err = g.client.Group.DeleteGroup(ctx, f.GroupID)
	if err != nil {
		return err
	}

	_, err = g.client.User.DeactivateUser(ctx, f.User, nil)
	if err != nil {
		return err
	}

	_, err = g.client.User.DeactivateOrDeleteUser(ctx, f.User, nil)
	return err
}
