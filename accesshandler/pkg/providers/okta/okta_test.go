package okta

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
	"github.com/hashicorp/go-multierror"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestIntegration(t *testing.T) {
	ctx := context.Background()
	_ = godotenv.Load("../../../.env")
	if os.Getenv("GRANTED_INTEGRATION_TEST") == "" {
		t.Skip("GRANTED_INTEGRATION_TEST is not set, skipping integration testing")
	}
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
		{
			Name:              "group not exist",
			Subject:           f.User,
			Args:              `{"groupId": "non-existent"}`,
			WantValidationErr: &multierror.Error{Errors: []error{&GroupNotFoundError{Group: "non-existent"}}},
		},
		{
			Name:              "subject not exist",
			Subject:           "other",
			Args:              fmt.Sprintf(`{"groupId": "%s"}`, f.GroupID),
			WantValidationErr: &multierror.Error{Errors: []error{&UserNotFoundError{User: "other"}}},
		},
		{
			Name:              "group and subject not exist",
			Subject:           "other",
			Args:              `{"groupId": "non-existent"}`,
			WantValidationErr: &multierror.Error{Errors: []error{&UserNotFoundError{User: "other"}, &GroupNotFoundError{Group: "non-existent"}}},
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
