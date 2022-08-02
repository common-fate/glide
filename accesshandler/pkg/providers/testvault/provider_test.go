package testvault

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers/okta/fixtures"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providertest"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providertest/integration"
	"github.com/common-fate/granted-approvals/pkg/gconfig"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestIntegration(t *testing.T) {
	if os.Getenv("GRANTED_INTEGRATION_TEST") == "" {
		t.Skip("GRANTED_INTEGRATION_TEST is not set, skipping integration testing")
	}

	ctx := context.Background()
	_ = godotenv.Load("../../../.env")

	var f fixtures.Fixtures
	err := providertest.LoadFixture(ctx, "okta", &f)
	if err != nil {
		t.Fatal(err)
	}

	testcases := []integration.TestCase{
		{
			Name:              "ok",
			Subject:           f.User,
			Args:              fmt.Sprintf(`{"groupId": "%s"}`, f.GroupID),
			WantValidationErr: nil,
		},
	}
	pc := os.Getenv("PROVIDER_CONFIG")
	var configMap map[string]json.RawMessage
	err = json.Unmarshal([]byte(pc), &configMap)
	if err != nil {
		t.Fatal(err)
	}
	integration.RunTests(t, ctx, "okta", &Provider{}, testcases, integration.WithProviderConfig(configMap["okta"]))
}

func TestArgSchema(t *testing.T) {
	o := Provider{}

	res := o.ArgSchema()
	out, err := json.Marshal(res)
	if err != nil {
		t.Fatal(err)
	}
	want, err := ioutil.ReadFile("./testdata/argschema.json")
	if err != nil {
		t.Fatal(err)
	}
	buffer := new(bytes.Buffer)
	err = json.Compact(buffer, want)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, buffer.String(), string(out))
}

func TestInstructions(t *testing.T) {
	p := Provider{
		apiURL:   gconfig.StringValue{Value: "https://testvault.internal"},
		uniqueID: gconfig.StringValue{Value: "1234"},
	}
	args := `{"vault": "my-vault"}`
	got, err := p.Instructions(context.Background(), "testuser", []byte(args))
	if err != nil {
		t.Fatal(err)
	}
	want := "This is just a test resource to show you how Granted Approvals works.\nVisit the [vault membership URL](https://testvault.internal/vaults/1234_my-vault/members/testuser) to check that your access has been provisioned."
	assert.Equal(t, want, got)
}
