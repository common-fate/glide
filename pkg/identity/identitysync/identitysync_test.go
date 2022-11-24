package identitysync

import (
	"context"
	"os"
	"testing"

	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/common-fate/granted-approvals/pkg/gconfig"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

type identitySyncTestCase struct {
	Name    string
	idpType string
}

func TestIntegration(t *testing.T) {
	ctx := context.Background()
	_ = godotenv.Load("../../../.env")

	if os.Getenv("COMMON_FATE_INTEGRATION_TEST") == "" {
		t.Skip("COMMON_FATE_INTEGRATION_TEST is not set, skipping integration testing")
	}
	idpConfig := os.Getenv("IDENTITY_SETTINGS")
	if idpConfig == "" {
		t.Skip("IDENTITY_SETTINGS is not set, skipping integration testing")
	}

	ic, err := deploy.UnmarshalFeatureMap(idpConfig)
	if err != nil {
		panic(err)
	}

	testcases := []identitySyncTestCase{
		{
			Name:    "OneLogin ok",
			idpType: "one-login",
		},
		{
			Name:    "Azure ok",
			idpType: "azure",
		},
		{
			Name:    "Okta ok",
			idpType: "okta",
		},
		{
			Name:    "Gsuite ok",
			idpType: "google",
		},
		{
			Name:    "AWS SSO ok",
			idpType: "aws-sso",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {

			idp, ok := Registry().IdentityProviders[tc.idpType]
			assert.True(t, ok)

			cfg := idp.IdentityProvider.Config()
			idpCfg, ok := ic[tc.idpType]
			if !ok {
				t.Skip("identity config for idp type is missing", tc.idpType)
			}

			//set the config
			err := cfg.Load(ctx, &gconfig.MapLoader{Values: idpCfg})
			assert.NoError(t, err)

			err = idp.IdentityProvider.Init(ctx)

			assert.NoError(t, err)

			if tester, ok := idp.IdentityProvider.(gconfig.Tester); ok {
				err = tester.TestConfig(ctx)
				assert.NoError(t, err)
			}
			grps, err := idp.IdentityProvider.ListGroups(ctx)
			assert.NoError(t, err)
			assert.Greater(t, len(grps), 0)

			usrs, err := idp.IdentityProvider.ListUsers(ctx)
			assert.NoError(t, err)
			assert.Greater(t, len(usrs), 0)

		})
	}
}
