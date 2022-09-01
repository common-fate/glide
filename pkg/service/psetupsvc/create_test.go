package psetupsvc

import (
	"context"
	"testing"

	"github.com/common-fate/ddb/ddbmock"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providerregistry"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/psetup"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/common-fate/granted-approvals/pkg/providersetup"
	"github.com/common-fate/granted-approvals/pkg/storage"
	"github.com/common-fate/granted-approvals/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	type testcase struct {
		name                   string
		providerType           string
		existingProviders      deploy.ProviderMap
		templateData           psetup.TemplateData
		registry               providerregistry.ProviderRegistry
		existingProviderSetups []providersetup.Setup
		want                   *providersetup.Setup
		wantErr                error
	}

	testcases := []testcase{
		{
			name:         "ok",
			providerType: "commonfate/test",
			registry: providerregistry.ProviderRegistry{
				Providers: map[string]map[string]providerregistry.RegisteredProvider{
					"commonfate/test": {
						"v1": {
							DefaultID: "test",
						},
					},
				},
			},
			want: &providersetup.Setup{
				ID:               "test",
				Status:           types.INITIALCONFIGURATIONINPROGRESS,
				ProviderType:     "commonfate/test",
				ProviderVersion:  "v1",
				ConfigValues:     map[string]string{},
				ConfigValidation: map[string]providersetup.Validation{},
			},
		},
		{
			name:         "increment ID for new provider",
			providerType: "commonfate/test",
			registry: providerregistry.ProviderRegistry{
				Providers: map[string]map[string]providerregistry.RegisteredProvider{
					"commonfate/test": {
						"v1": {
							DefaultID: "test",
						},
					},
				},
			},
			existingProviders: deploy.ProviderMap{
				"test": deploy.Provider{
					Uses: "commonfate/test@v1",
				},
			},
			want: &providersetup.Setup{
				ID:               "test-2",
				Status:           types.INITIALCONFIGURATIONINPROGRESS,
				ProviderType:     "commonfate/test",
				ProviderVersion:  "v1",
				ConfigValues:     map[string]string{},
				ConfigValidation: map[string]providersetup.Validation{},
			},
		},
		{
			name:         "increment ID for new provider with pending setup",
			providerType: "commonfate/test",
			registry: providerregistry.ProviderRegistry{
				Providers: map[string]map[string]providerregistry.RegisteredProvider{
					"commonfate/test": {
						"v1": {
							DefaultID: "test",
						},
					},
				},
			},
			// we have an existing provider registered in the config as 'test'
			existingProviders: deploy.ProviderMap{
				"test": deploy.Provider{
					Uses: "commonfate/test@v1",
				},
			},
			// we've also got a pending setup for 'test-2'
			existingProviderSetups: []providersetup.Setup{
				{ID: "test-2", ProviderType: "commonfate/test", ProviderVersion: "v1"},
			},
			want: &providersetup.Setup{
				// should be 'test-3' as both 'test' and 'test-2' are taken
				ID:               "test-3",
				Status:           types.INITIALCONFIGURATIONINPROGRESS,
				ProviderType:     "commonfate/test",
				ProviderVersion:  "v1",
				ConfigValues:     map[string]string{},
				ConfigValidation: map[string]providersetup.Validation{},
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			db := ddbmock.New(t)
			db.MockQuery(&storage.ListProviderSetupsForType{Result: tc.existingProviderSetups})

			s := Service{
				DB:           db,
				TemplateData: tc.templateData,
			}
			got, err := s.Create(ctx, tc.providerType, tc.existingProviders, tc.registry)
			assert.Equal(t, tc.want, got)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
