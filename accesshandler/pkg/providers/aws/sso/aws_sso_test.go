package sso

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers/aws/sso/fixtures"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providertest"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providertest/integration"
	"github.com/joho/godotenv"
)

func TestIntegration(t *testing.T) {
	ctx := context.Background()
	_ = godotenv.Load("../../../../.env")
	if os.Getenv("GRANTED_INTEGRATION_TEST") == "" {
		t.Skip("GRANTED_INTEGRATION_TEST is not set, skipping integration testing")
	}

	var f fixtures.Fixtures
	err := providertest.LoadFixture(ctx, "aws_sso", &f)
	if err != nil {
		t.Fatal(err)
	}

	testcases := []integration.TestCase{
		{
			Name:              "ok",
			Subject:           f.User,
			Args:              fmt.Sprintf(`{"permissionSetArn": "%s", "accountId": "%s"}`, f.PermissionSetARN, f.AccountID),
			WantValidationErr: nil,
			WantValidationSucceeded: map[string]bool{
				"user-exists-in-okta":  true,
				"group-exists-in-okta": false,
			},
		},
		{
			Name:              "permissionSet not exist",
			Subject:           f.User,
			Args:              fmt.Sprintf(`{"permissionSetArn": "arn:aws:sso:::permissionSet/ssoins-1234561234512345/ps-1234512345123451", "accountId": "%s"}`, f.AccountID),
			WantValidationErr: &PermissionSetNotFoundErr{PermissionSet: "arn:aws:sso:::permissionSet/ssoins-1234561234512345/ps-1234512345123451"},
			WantValidationSucceeded: map[string]bool{
				"user-exists-in-okta":  true,
				"group-exists-in-okta": false,
			},
		},
		{
			Name:              "subject not exist",
			Subject:           "other",
			Args:              fmt.Sprintf(`{"permissionSetArn": "%s", "accountId": "%s"}`, f.PermissionSetARN, f.AccountID),
			WantValidationErr: &UserNotFoundError{Email: "other"},
			WantValidationSucceeded: map[string]bool{
				"user-exists-in-okta":  true,
				"group-exists-in-okta": false,
			},
		},
		{
			Name:              "account not exist",
			Subject:           f.User,
			Args:              fmt.Sprintf(`{"permissionSetArn": "%s", "accountId": "123456789012"}`, f.PermissionSetARN),
			WantValidationErr: &AccountNotFoundError{AccountID: "123456789012"},
			WantValidationSucceeded: map[string]bool{
				"user-exists-in-okta":  true,
				"group-exists-in-okta": false,
			},
		},
	}
	pc := os.Getenv("PROVIDER_CONFIG")
	var configMap map[string]json.RawMessage
	err = json.Unmarshal([]byte(pc), &configMap)
	if err != nil {
		t.Fatal(err)
	}
	integration.RunTests(t, ctx, "aws_sso", &Provider{}, testcases, integration.WithProviderConfig(configMap["aws_sso"]))
}
