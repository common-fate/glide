package identitysync

import (
	"context"
	"testing"

	"github.com/common-fate/granted-approvals/pkg/gconfig"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

type idpTestCase struct {
	Name    string
	idpType string
	config  map[string]string
}

func TestIntegration(t *testing.T) {
	ctx := context.Background()
	err := godotenv.Load("../../../.env")
	assert.NoError(t, err)

	testcases := []idpTestCase{
		{
			Name:    "list users ok",
			idpType: "one-login",
			config: map[string]string{
				"baseURL":      "https://commonfate-dev.onelogin.com",
				"clientId":     "ec6dad650566db9f4f12241f9b55ad18220be45c7049feb8b667db16cc01f36e",
				"clientSecret": "awsssm:///granted/secrets/identity/one-login/secret:6",
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {

			idp, ok := Registry().IdentityProviders[tc.idpType]
			assert.True(t, ok)

			cfg := idp.IdentityProvider.Config()

			//set the config
			err := cfg.Load(ctx, &gconfig.MapLoader{Values: tc.config})
			assert.NoError(t, err)

			err = idp.IdentityProvider.Init(ctx)

			assert.NoError(t, err)

			grps, err := idp.IdentityProvider.ListGroups(ctx)
			assert.NoError(t, err)
			assert.Greater(t, len(grps), 0)

		})
	}
}
