package ssov2

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/common-fate/common-fate/accesshandler/pkg/providers/okta/fixtures"
	"github.com/common-fate/common-fate/accesshandler/pkg/providertest"
	"github.com/common-fate/common-fate/accesshandler/pkg/providertest/integration"
	"github.com/joho/godotenv"
)

func TestIntegration(t *testing.T) {
	ctx := context.Background()
	_ = godotenv.Load("../../../../../.env")
	if os.Getenv("COMMONFATE_INTEGRATION_TEST") == "" {
		t.Skip("COMMONFATE_INTEGRATION_TEST is not set, skipping integration testing")
	}
	var f fixtures.Fixtures
	err := providertest.LoadFixture(ctx, "aws-sso-v2", &f)
	if err != nil {
		t.Fatal(err)
	}

	testcases := []integration.TestCase{
		{
			Name:              "ok",
			Subject:           f.User,
			Args:              fmt.Sprintf(`{"accountId": "%s", "permissionSetArn": "%s"}`, f.AccountID, f.PermissionSetARN),
			WantValidationErr: nil,
		},
		{
			Name:              "account not exist",
			Subject:           f.User,
			Args:              fmt.Sprintf(`{"accountId": "non-existent", "permissionSetArn": "%s"}`, f.PermissionSetARN),
			WantValidationErr: errors.New(""),
		},
		{
			Name:              "permission set not exist",
			Subject:           f.User,
			Args:              fmt.Sprintf(`{"accountId": "%s", "permissionSetArn": "non-existent"}`, f.AccountID),
			WantValidationErr: errors.New(""),
		},
		{
			Name:              "subject not exist",
			Subject:           "other",
			Args:              fmt.Sprintf(`{"accountId": "%s", "permissionSetArn": "%s"}`, f.AccountID, f.PermissionSetARN),
			WantValidationErr: errors.New(""),
		},
		{
			Name:              "permission set and subject not exist",
			Subject:           "other",
			Args:              fmt.Sprintf(`{"accountId": "%s", "permissionSetArn": "non-existent"}`, f.AccountID),
			WantValidationErr: errors.New(""),
		},
	}
	w, err := integration.ProviderWith("aws-sso-v2")
	if err != nil {
		t.Fatal(err)
	}
	integration.RunTests(t, ctx, "aws-sso-v2", &Provider{}, testcases, integration.WithProviderConfig(w))
}
