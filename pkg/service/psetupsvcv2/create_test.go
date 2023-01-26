package psetupsvcv2

import (
	"testing"

	"github.com/common-fate/common-fate/accesshandler/pkg/providerregistry"
	"github.com/common-fate/common-fate/pkg/deploy"
	"github.com/common-fate/common-fate/pkg/providersetupv2"
	"github.com/common-fate/common-fate/pkg/types"
)

func TestCreate(t *testing.T) {
	type testcase struct {
		name              string
		providerType      string
		existingProviders deploy.ProviderMap
		// templateData           psetup.TemplateData
		registry               providerregistry.ProviderRegistry
		existingProviderSetups []providersetupv2.Setup
		want                   *providersetupv2.Setup
		// wantErr                error
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
			want: &providersetupv2.Setup{
				ID:               "test",
				Status:           types.ProviderSetupV2StatusINITIALCONFIGURATIONINPROGRESS,
				ProviderTeam:     "commonfate",
				ProviderName:     "test",
				ProviderVersion:  "v1",
				ConfigValues:     map[string]string{},
				ConfigValidation: map[string]providersetupv2.Validation{},
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
			want: &providersetupv2.Setup{
				ID:               "test-2",
				Status:           types.ProviderSetupV2StatusINITIALCONFIGURATIONINPROGRESS,
				ProviderTeam:     "commonfate",
				ProviderName:     "test",
				ProviderVersion:  "v1",
				ConfigValues:     map[string]string{},
				ConfigValidation: map[string]providersetupv2.Validation{},
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
			existingProviderSetups: []providersetupv2.Setup{
				{ID: "test-2",
					ProviderTeam: "commonfate",
					ProviderName: "test", ProviderVersion: "v1"},
			},
			want: &providersetupv2.Setup{
				// should be 'test-3' as both 'test' and 'test-2' are taken
				ID:               "test-3",
				Status:           types.ProviderSetupV2StatusINITIALCONFIGURATIONINPROGRESS,
				ProviderTeam:     "commonfate",
				ProviderName:     "test",
				ProviderVersion:  "v1",
				ConfigValues:     map[string]string{},
				ConfigValidation: map[string]providersetupv2.Validation{},
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// ctx := context.Background()
			// db := ddbmock.New(t)
			// db.MockQuery(&storage.ListProviderSetupsForType{Result: tc.existingProviderSetups})

			// s := Service{
			// 	DB:           db,
			// 	TemplateData: tc.templateData,
			// }
			// got, err := s.Create(ctx, tc.providerType, tc.existingProviders, tc.registry)
			// assert.Equal(t, tc.want, got)
			// assert.Equal(t, tc.wantErr, err)
		})
	}
}
