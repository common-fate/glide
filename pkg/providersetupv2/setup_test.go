package providersetupv2

import (
	"testing"

	ahtypes "github.com/common-fate/common-fate/accesshandler/pkg/types"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestUpdateValidationStatus(t *testing.T) {
	type testcase struct {
		name string
		give Setup
		want types.ProviderSetupV2Status
	}

	testcases := []testcase{
		{
			name: "error",
			give: Setup{
				Status: types.ProviderSetupV2StatusINITIALCONFIGURATIONINPROGRESS,
				ConfigValidation: map[string]Validation{
					"test": {
						Status: ahtypes.ERROR,
					},
				},
			},
			want: types.ProviderSetupV2StatusVALIDATIONFAILED,
		},
		{
			name: "success",
			give: Setup{
				Status: types.ProviderSetupV2StatusINITIALCONFIGURATIONINPROGRESS,
				ConfigValidation: map[string]Validation{
					"test": {
						Status: ahtypes.SUCCESS,
					},
				},
			},
			want: types.ProviderSetupV2StatusVALIDATIONSUCEEDED,
		},
		{
			name: "other",
			give: Setup{
				Status: types.ProviderSetupV2StatusINITIALCONFIGURATIONINPROGRESS,
				ConfigValidation: map[string]Validation{
					"test": {
						Status: ahtypes.INPROGRESS,
					},
				},
			},
			want: types.ProviderSetupV2StatusINITIALCONFIGURATIONINPROGRESS,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			setup := tc.give
			setup.UpdateValidationStatus()
			assert.Equal(t, tc.want, setup.Status)
		})
	}
}
