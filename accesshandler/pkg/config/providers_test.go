package config

import (
	"context"
	"errors"
	"testing"

	"github.com/common-fate/granted-approvals/accesshandler/pkg/genv"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers/aws/sso"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers/okta"
	"github.com/stretchr/testify/assert"
)

// testProvider configures a provider with testing variables based on the provided 'vals' argument.
func testProvider(t *testing.T, p providers.Accessor, vals map[string]string) providers.Accessor {
	ctx := context.Background()

	if c, ok := p.(providers.Configer); ok {
		err := c.Config().Load(ctx, &genv.MapLoader{
			Values: vals,
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	// initialise the provider if it supports it.
	if i, ok := p.(providers.Initer); ok {
		err := i.Init(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}

	return p
}

func TestConfigureProviders(t *testing.T) {
	ctx := context.Background()

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
			give: `{"test": {"uses": "commonfate/aws-sso@v1", "with": {"identityStoreId": "id-123", "instanceArn": "arn::test", "region": "us-east-1"}}}`,
			want: map[string]Provider{
				"test": {
					ID:   "test",
					Type: "aws-sso",
					Provider: testProvider(t, &sso.Provider{}, map[string]string{
						"identityStoreId": "id-123",
						"instanceArn":     "arn::test",
						"region":          "us-east-1",
					}),
				},
			},
		},
		{
			name: "aws sso with no region",
			give: `{"test": {"uses": "commonfate/aws-sso@v1", "with": {"identityStoreId": "id-123", "instanceArn": "arn::test"}}}`,
			want: map[string]Provider{
				"test": {
					ID:   "test",
					Type: "aws-sso",
					Provider: testProvider(t, &sso.Provider{}, map[string]string{
						"identityStoreId": "id-123",
						"instanceArn":     "arn::test",
						"region":          "",
					}),
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := ConfigureProviders(ctx, []byte(tc.give))
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

				if c, ok := p.Provider.(providers.Configer); ok {
					gotc := got.Provider.(providers.Configer)
					assert.Equal(t, c.Config(), gotc.Config())
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
