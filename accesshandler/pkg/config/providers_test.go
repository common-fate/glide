package config

import (
	"context"
	"errors"
	"testing"

	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers"
	ssov2 "github.com/common-fate/granted-approvals/accesshandler/pkg/providers/aws/sso-v2"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers/okta"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/common-fate/granted-approvals/pkg/gconfig"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

// testProvider configures a provider with testing variables based on the provided 'vals' argument.
func testProvider(t *testing.T, p providers.Accessor, vals map[string]string) providers.Accessor {
	ctx := context.Background()

	if c, ok := p.(gconfig.Configer); ok {
		err := c.Config().Load(ctx, &gconfig.MapLoader{
			Values: vals,
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	// initialise the provider if it supports it.
	if i, ok := p.(gconfig.Initer); ok {
		err := i.Init(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}

	return p
}

// Note that this test requires AWS credentials to be in the environment
func TestConfigureProviders(t *testing.T) {
	ctx := context.Background()
	_ = godotenv.Load("../../../.env")
	type testcase struct {
		name string
		give string
		want map[string]Provider
	}

	testcases := []testcase{
		{
			name: "okta",
			give: `{"test": {"uses": "commonfate/okta@v1", "with": {"orgUrl": "https://test.internal", "apiToken": "secret"}}}`,
			want: map[string]Provider{
				"test": {
					ID:   "test",
					Type: "okta",
					Provider: testProvider(t, &okta.Provider{}, map[string]string{
						"orgUrl":   "https://test.internal",
						"apiToken": "secret",
					}),
				},
			},
		},
		{
			name: "aws sso",
			give: `{"test": {"uses": "commonfate/aws-sso@v2", "with": {"identityStoreId": "id-123", "instanceArn": "arn::test", "region": "us-east-1", "ssoRoleArn": "arn::test"}}}`,
			want: map[string]Provider{
				"test": {
					ID:      "test",
					Type:    "aws-sso",
					Version: "v2",

					Provider: testProvider(t, &ssov2.Provider{}, map[string]string{
						"identityStoreId": "id-123",
						"instanceArn":     "arn::test",
						"region":          "us-east-1",
						"ssoRoleArn":      "arn::test",
					}),
				},
			},
		},
		{
			name: "aws sso with no region",
			give: `{"test": {"uses": "commonfate/aws-sso@v2", "with": {"identityStoreId": "id-123", "instanceArn": "arn::test", "ssoRoleArn": "arn::test"}}}`,
			want: map[string]Provider{
				"test": {
					ID:      "test",
					Type:    "aws-sso",
					Version: "v2",
					Provider: testProvider(t, &ssov2.Provider{}, map[string]string{
						"identityStoreId": "id-123",
						"instanceArn":     "arn::test",
						"region":          "",
						"ssoRoleArn":      "arn::test",
					}),
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			cfg, err := deploy.UnmarshalProviderMap(tc.give)
			if err != nil {
				t.Fatal(err)
			}

			err = ConfigureProviders(ctx, cfg)
			if err != nil {
				t.Fatal(err)
			}
			for k, p := range tc.want {
				got, ok := Providers[k]
				if !ok {
					t.Fatalf("did not load provider %s", k)
				}
				assert.Equal(t, p.ID, got.ID)
				assert.Equal(t, p.Type, got.Type)
				assert.IsType(t, p.Provider, got.Provider)

				if c, ok := p.Provider.(gconfig.Configer); ok {
					gotc := got.Provider.(gconfig.Configer)
					assert.Len(t, gotc.Config(), len(c.Config()))
					for _, v := range c.Config() {
						found := false
						for _, v1 := range gotc.Config() {
							if v.Key() == v1.Key() {
								found = true
								assert.Equal(t, v.Get(), v1.Get())
							}
						}
						assert.True(t, found)

					}
				}
			}
		})
	}
}

func TestProviderFromUses(t *testing.T) {
	type testcase struct {
		name    string
		give    string
		want    Provider
		wantErr error
	}

	testcases := []testcase{
		{
			name: "ok",
			give: "commonfate/test@v1",
			want: Provider{
				Type:    "test",
				Version: "v1",
			},
		},
		{
			name:    "bad input",
			give:    "commonfate",
			wantErr: errors.New("could not extract provider information from commonfate"),
		},
		{
			name:    "no version",
			give:    "commonfate/test",
			wantErr: errors.New("could not extract provider information from commonfate/test"),
		},
		{
			name:    "no version with @",
			give:    "commonfate/test@",
			wantErr: errors.New("could not extract provider version from commonfate/test@"),
		},
		{
			name: "with special characters",
			give: "testing_-test/test-_provider@v1.1",
			want: Provider{
				Type:    "test-_provider",
				Version: "v1.1",
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := providerFromUses(tc.give)
			if err != nil && tc.wantErr == nil {
				t.Fatal(err)
			}
			if tc.wantErr != nil {
				assert.Equal(t, tc.wantErr, err)
			}
			assert.Equal(t, tc.want, got)
		})
	}

}
