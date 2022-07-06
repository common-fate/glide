package sso

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers/aws/sso/fixtures"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providertest"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providertest/integration"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestIntegration(t *testing.T) {
	if os.Getenv("GRANTED_INTEGRATION_TEST") == "" {
		t.Skip("GRANTED_INTEGRATION_TEST is not set, skipping integration testing")
	}

	ctx := context.Background()
	_ = godotenv.Load("../../../../.env")

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
		},
		{
			Name:              "permissionSet not exist",
			Subject:           f.User,
			Args:              fmt.Sprintf(`{"permissionSetArn": "arn:aws:sso:::permissionSet/ssoins-1234561234512345/ps-1234512345123451", "accountId": "%s"}`, f.AccountID),
			WantValidationErr: &PermissionSetNotFoundErr{PermissionSet: "arn:aws:sso:::permissionSet/ssoins-1234561234512345/ps-1234512345123451"},
		},
		{
			Name:              "subject not exist",
			Subject:           "other",
			Args:              fmt.Sprintf(`{"permissionSetArn": "%s", "accountId": "%s"}`, f.PermissionSetARN, f.AccountID),
			WantValidationErr: &UserNotFoundError{Email: "other"},
		},
		{
			Name:              "account not exist",
			Subject:           f.User,
			Args:              fmt.Sprintf(`{"permissionSetArn": "%s", "accountId": "123456789012"}`, f.PermissionSetARN),
			WantValidationErr: &AccountNotFoundError{AccountID: "123456789012"},
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

func TestArgSchema(t *testing.T) {
	p := Provider{}

	res := p.ArgSchema()
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
