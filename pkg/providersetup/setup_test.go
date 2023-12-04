package providersetup

import (
	"testing"

	ahtypes "github.com/common-fate/common-fate/accesshandler/pkg/types"
	"github.com/common-fate/common-fate/pkg/deploy"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestUpdateValidationStatus(t *testing.T) {
	type testcase struct {
		name string
		give Setup
		want types.ProviderSetupStatus
	}

	testcases := []testcase{
		{
			name: "error",
			give: Setup{
				Status: types.INITIALCONFIGURATIONINPROGRESS,
				ConfigValidation: map[string]Validation{
					"test": {
						Status: ahtypes.ERROR,
					},
				},
			},
			want: types.VALIDATIONFAILED,
		},
		{
			name: "success",
			give: Setup{
				Status: types.INITIALCONFIGURATIONINPROGRESS,
				ConfigValidation: map[string]Validation{
					"test": {
						Status: ahtypes.SUCCESS,
					},
				},
			},
			want: types.VALIDATIONSUCEEDED,
		},
		{
			name: "other",
			give: Setup{
				Status: types.INITIALCONFIGURATIONINPROGRESS,
				ConfigValidation: map[string]Validation{
					"test": {
						Status: ahtypes.INPROGRESS,
					},
				},
			},
			want: types.INITIALCONFIGURATIONINPROGRESS,
		},
	}

	for i := range testcases {
		tc := testcases[i]
		t.Run(tc.name, func(t *testing.T) {
			setup := tc.give
			setup.UpdateValidationStatus()
			assert.Equal(t, tc.want, setup.Status)
		})
	}
}

func TestToProvider(t *testing.T) {
	type testcase struct {
		name string
		give Setup
		want deploy.Provider
	}

	testcases := []testcase{
		{
			name: "ok",
			give: Setup{
				ProviderType:    "test",
				ProviderVersion: "v1",
				ConfigValues: map[string]string{
					"value": "testing",
				},
			},
			want: deploy.Provider{
				Uses: "test@v1",
				With: map[string]string{
					"value": "testing",
				},
			},
		},
	}

	for i := range testcases {
		tc := testcases[i]
		t.Run(tc.name, func(t *testing.T) {
			got := tc.give.ToProvider()
			assert.Equal(t, tc.want, got)
		})
	}
}
