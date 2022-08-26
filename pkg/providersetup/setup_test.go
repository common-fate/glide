package providersetup

import (
	"testing"

	ahtypes "github.com/common-fate/granted-approvals/accesshandler/pkg/types"
	"github.com/common-fate/granted-approvals/pkg/types"
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

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			setup := tc.give
			setup.UpdateValidationStatus()
			assert.Equal(t, tc.want, setup.Status)
		})
	}
}
