package ad

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers/okta/fixtures"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providertest"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providertest/integration"
	"github.com/hashicorp/go-multierror"
	"github.com/joho/godotenv"
)

func TestIntegration(t *testing.T) {
	if os.Getenv("GRANTED_INTEGRATION_TEST") == "" {
		t.Skip("GRANTED_INTEGRATION_TEST is not set, skipping integration testing")
	}

	ctx := context.Background()
	_ = godotenv.Load("../../../../.env")

	var f fixtures.Fixtures
	err := providertest.LoadFixture(ctx, "azure", &f)
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
	integration.RunTests(t, ctx, "azure", &Provider{}, testcases, integration.WithProviderConfig(configMap["azure"]))
}
