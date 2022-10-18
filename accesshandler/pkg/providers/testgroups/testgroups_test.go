package testgroups

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/common-fate/granted-approvals/accesshandler/pkg/providertest/integration"
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
