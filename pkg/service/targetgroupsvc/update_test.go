package targetgroupsvc

import (
	"context"
	"net/http"
	"testing"

	"github.com/benbjohnson/clock"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/targetgroup"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb/ddbmock"
	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk/prmocks"
	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
)

func TestUpdateTargetGroup(t *testing.T) {
	type testcase struct {
		name         string
		version      string
		give         UpdateOpts
		wantErr      error
		withResponse providerregistrysdk.Provider
		want         *targetgroup.TargetGroup
		tgLookup     targetgroup.TargetGroup
	}

	mockTargetGroup := targetgroup.TargetGroup{
		ID:           "test",
		TargetSchema: targetgroup.GroupTargetSchema{From: "commonfate/test/v1.0.2", Schema: providerregistrysdk.TargetSchema{}},
	}

	mockTargetGroupBefore := targetgroup.TargetGroup{
		ID:           "test",
		TargetSchema: targetgroup.GroupTargetSchema{From: "commonfate/test/v1.0.1", Schema: providerregistrysdk.TargetSchema{}},
	}

	testcases := []testcase{
		{
			name:    "ok",
			version: "v1.0.2",
			give:    UpdateOpts{UpdateRequest: types.CreateTargetGroupRequest{ID: "test", TargetSchema: "commonfate/test/v1.0.2"}},
			withResponse: providerregistrysdk.Provider{
				Team:    "commonfate",
				Name:    "test",
				Version: "v1.0.2",
				Schema:  providerregistrysdk.ProviderSchema{Target: providerregistrysdk.TargetSchema{}},
			},
			want:     &mockTargetGroup,
			tgLookup: mockTargetGroupBefore,
		},
		// {
		// 	name:    "provider does not exist",
		// 	version: "v2.2.2",

		// 	give: UpdateOpts{UpdateRequest: types.CreateTargetGroupRequest{ID: "test", TargetSchema: "commonfate/test/v1.0.2"}},

		// 	wantErr: apio.NewRequestError(errors.New("provider does not exist"), http.StatusBadRequest),
		// 	want:    nil,
		// },
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			db := ddbmock.New(t)
			db.MockQuery(&storage.GetTargetGroup{Result: tc.tgLookup})

			clk := clock.NewMock()
			ctrl := gomock.NewController(t)

			defer ctrl.Finish()

			m := prmocks.NewMockClientWithResponsesInterface(ctrl)
			m.EXPECT().GetProviderWithResponse(gomock.Any(), gomock.Eq("commonfate"), gomock.Eq("test"), gomock.Eq(tc.version)).Return(&providerregistrysdk.GetProviderResponse{HTTPResponse: &http.Response{StatusCode: 200}, JSON200: &tc.withResponse}, nil)

			s := Service{
				Clock:                  clk,
				DB:                     db,
				ProviderRegistryClient: m,
			}

			got, err := s.UpdateTargetGroup(context.Background(), tc.give)

			if err != nil {
				assert.Equal(t, tc.wantErr.Error(), err.Error())
			}
			assert.Equal(t, tc.want, got)

		})
	}

}
