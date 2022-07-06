package testgroups

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/common-fate/granted-approvals/accesshandler/pkg/providertest/integration"
	"github.com/stretchr/testify/assert"
)

func TestIntegration(t *testing.T) {
	if os.Getenv("GRANTED_INTEGRATION_TEST") == "" {
		t.Skip("GRANTED_INTEGRATION_TEST is not set, skipping integration testing")
	}

	ctx := context.Background()

	b, err := json.Marshal(map[string]interface{}{
		"groups": []string{"group1", "group2", "group3"},
	})
	if err != nil {
		t.Fatal(err)
	}
	testcases := []integration.TestCase{
		{
			Name:              "ok",
			Subject:           "test",
			Args:              `{"group": "group1"}`,
			WantValidationErr: nil,
		},
		{
			Name:              "group not exist",
			Subject:           "test",
			Args:              `{"group": "non-existent"}`,
			WantValidationErr: &GroupNotFoundError{Group: "non-existent"},
		},
	}

	integration.RunTests(t, ctx, "testgroups", &Provider{}, testcases, integration.WithProviderConfig(b))
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
