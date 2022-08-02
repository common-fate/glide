package lambda

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/common-fate/granted-approvals/accesshandler/pkg/config"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/types"
	"github.com/common-fate/granted-approvals/pkg/gconfig"
	"github.com/common-fate/iso8601"
	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/joho/godotenv"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
	"github.com/stretchr/testify/assert"
)

type TestConfig struct {
	ProviderID string
	GroupID    string
	Email      string
	OrgURL     string
	APIToken   string
}

func TestRevokeGrant(t *testing.T) {
	ctx := context.Background()

	_ = godotenv.Load("../../../../.env")
	testCfg := TestConfig{}

	loader := gconfig.EnvLoader{Prefix: "REVOKE_GRANT_INTEGRATION_TEST_"}

	err := loader.Load(
		gconfig.String("PROVIDER_ID", &testCfg.ProviderID, ""),
		gconfig.String("GROUP_ID", &testCfg.GroupID, ""),
		gconfig.String("SUBJECT_EMAIL", &testCfg.Email, ""),
		gconfig.String("OKTA_ORG_URL", &testCfg.OrgURL, ""),
		gconfig.String("OKTA_SYNC_TOKEN", &testCfg.APIToken, ""))
	if err != nil {
		t.Skip("environment variables not set")
	}

	//create a new grant using create grant
	runtime := Runtime{}

	err = runtime.Init(ctx)
	if err != nil {
		t.Fatal(err)
	}

	cfg, err := config.ReadProviderConfig(ctx, "lambda")
	if err != nil {
		t.Fatal(err)
	}

	err = config.ConfigureProviders(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}

	_, client, err := okta.NewClient(ctx, okta.WithOrgUrl(testCfg.OrgURL), okta.WithToken(testCfg.APIToken), okta.WithCache(false))
	if err != nil {
		t.Fatal(err)
	}
	//check that the user is not assigned to the group in okta

	users, res, err := client.Group.ListGroupUsers(ctx, testCfg.GroupID, query.NewQueryParams())

	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		assert.Equal(t, []*okta.User{}, users)
	}

	grant, err := runtime.CreateGrant(ctx, types.ValidCreateGrant{CreateGrant: types.CreateGrant{
		Id:       "TESTGRANT",
		Start:    iso8601.Now(),
		End:      iso8601.New(time.Now().Add(time.Minute)),
		Subject:  openapi_types.Email(testCfg.Email),
		Provider: testCfg.ProviderID,
		With:     types.CreateGrant_With{AdditionalProperties: map[string]string{"groupId": testCfg.GroupID}},
	},
	})

	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second * 10)

	//check that user is assigned to the group
	users, _, err = client.Group.ListGroupUsers(ctx, testCfg.GroupID, query.NewQueryParams())

	if err != nil {
		t.Fatal(err)
	}

	var userEmails []string
	for _, g := range users {
		userEmails = append(userEmails, (*g.Profile)["email"].(string))
	}

	assert.Contains(t, userEmails, testCfg.Email)

	//check the state function is running

	_, err = runtime.RevokeGrant(ctx, grant.ID, "actor")
	if err != nil {
		t.Fatal(err)
	}

	//check the okta group is not in the group
	users, _, err = client.Group.ListGroupUsers(ctx, testCfg.GroupID, query.NewQueryParams())

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, []*okta.User{}, users)

}
