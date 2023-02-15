package targetgroupsvc

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/benbjohnson/clock"
	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/common-fate/pkg/targetgroup"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb/ddbmock"
	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk/prmocks"
	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
)

func TestCreateTargetGroup(t *testing.T) {
	type testcase struct {
		name         string
		version      string
		give         types.CreateTargetGroupRequest
		wantErr      error
		withResponse providerregistrysdk.ProviderDetail
		want         *targetgroup.TargetGroup
	}

	mockTargetGroup := targetgroup.TargetGroup{
		ID:           "test",
		TargetSchema: targetgroup.GroupTargetSchema{From: "commonfate/test/v1.0.1", Schema: providerregistrysdk.TargetMode_Schema{}},
	}

	testcases := []testcase{
		{
			name:    "ok",
			version: "v1.0.1",
			give:    types.CreateTargetGroupRequest{ID: "test", TargetSchema: "commonfate/test/v1.0.1"},
			withResponse: providerregistrysdk.ProviderDetail{
				Publisher: "commonfate",
				Name:      "test",
				Version:   "v1.0.1",
			},
			want: &mockTargetGroup,
		},
		{
			name:    "provider does not exist",
			version: "v2.2.2",

			give: types.CreateTargetGroupRequest{ID: "test", TargetSchema: "commonfate/test/v2.2.2"},
			// withResponse: providerregistrysdk.Provider{
			// 	Team:    "commonfate",
			// 	Name:    "test",
			// 	Version: "v1.0.1",
			// 	Schema:  providerregistrysdk.ProviderSchema{},
			// },
			wantErr: apio.NewRequestError(errors.New("provider does not exist"), http.StatusBadRequest),
			want:    nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {

			dbc := ddbmock.Client{
				PutErr: tc.wantErr,
			}

			clk := clock.NewMock()
			ctrl := gomock.NewController(t)

			defer ctrl.Finish()

			m := prmocks.NewMockClientWithResponsesInterface(ctrl)
			m.EXPECT().GetProviderWithResponse(gomock.Any(), gomock.Eq("commonfate"), gomock.Eq("test"), gomock.Eq(tc.version)).Return(&providerregistrysdk.GetProviderResponse{HTTPResponse: &http.Response{StatusCode: 200}, JSON200: &tc.withResponse}, nil)

			s := Service{
				Clock:                  clk,
				DB:                     &dbc,
				ProviderRegistryClient: m,
			}

			got, err := s.CreateTargetGroup(context.Background(), tc.give)

			if err != nil {
				assert.Equal(t, tc.wantErr.Error(), err.Error())
			}
			assert.Equal(t, tc.want, got)

		})
	}

}
