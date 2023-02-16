package targetgroupsvc

import (
	"context"
	"testing"

	"github.com/benbjohnson/clock"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/targetgroup"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
	"github.com/common-fate/ddb/ddbmock"
	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk/prmocks"
	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
)

func TestCreateTargetGroup(t *testing.T) {
	// test cases:
	// s.DB.Query error == ddb.ErrNoItems -> ok ✅
	// s.DB.Query error == misc error     -> nil, misc error
	// s.DB.Query error == nil            -> nil, ErrTargetGroupIdAlreadyExists
	// invalid provider string			  -> nil, ErrInvalidProviderString
	// ProviderRegistryClient.GetProviderWithResponse error != nil -> nil, err
	// s.DB.Put error != nil -> nil, err

	// items to mock: s.DB.Query storage.GetTargetGroup ✅
	// items to mock: GetProviderWithResponse ✅
	// items to mock: s.DB.Put targetgroup.TargetGroup ✅

	// items to assert: tc.want, tc.wantErr ✅

	type testcase struct {
		name                         string
		version                      string
		give                         types.CreateTargetGroupRequest
		withResponse                 providerregistrysdk.ProviderDetail
		want                         *targetgroup.TargetGroup
		wantErr                      error
		mockStorageGetTargetGroupErr error
		mockGetProviderRes           providerregistrysdk.GetProviderResponse
		mockGetProviderResErr        error
		mockStoragePutTargetGroupErr error
	}

	clk := clock.NewMock()
	now := clk.Now()

	mockWantedTargetGroup := targetgroup.TargetGroup{
		ID: "test",
		TargetSchema: targetgroup.GroupTargetSchema{From: "commonfate/test@v1.0.1",
			Schema: providerregistrysdk.TargetMode_Schema{}},
		CreatedAt: now,
		UpdatedAt: now,
	}

	testcases := []testcase{
		{
			name:                         "ok",
			version:                      "v1.0.1",
			give:                         types.CreateTargetGroupRequest{ID: "test", TargetSchema: "commonfate/test@v1.0.1"},
			mockStorageGetTargetGroupErr: ddb.ErrNoItems,
			withResponse: providerregistrysdk.ProviderDetail{
				Publisher: "commonfate",
				Name:      "test",
				Version:   "v1.0.1",
			},
			mockGetProviderRes: providerregistrysdk.GetProviderResponse{
				JSON200: &providerregistrysdk.ProviderDetail{
					Version: "v1.0.1",
					Schema:  providerregistrysdk.ProviderSchema{},
				},
			},
			want: &mockWantedTargetGroup,
		},
		// {
		// 	name:                         "s.DB.Query error == misc error",
		// 	version:                      "v1.0.1",
		// 	give:                         types.CreateTargetGroupRequest{ID: "test", TargetSchema: "commonfate/test@v1.0.1"},
		// 	mockStorageGetTargetGroupErr: errors.New("error"),
		// 	mockGetProviderRes:           providerregistrysdk.GetProviderResponse{},
		// 	mockGetProviderResErr:        nil,
		// 	withResponse: providerregistrysdk.ProviderDetail{
		// 		Publisher: "commonfate",
		// 		Name:      "test",
		// 		Version:   "v1.0.1",
		// 	},
		// 	want:    nil,
		// 	wantErr: errors.New("error"),
		// },
	}

	for _, tc := range testcases {

		tc := tc

		t.Run(tc.name, func(t *testing.T) {

			dbc := ddbmock.New(t)

			dbc.MockQueryWithErr(&storage.GetTargetGroup{}, tc.mockStorageGetTargetGroupErr)

			ctrl := gomock.NewController(t)

			defer ctrl.Finish()

			m := prmocks.NewMockClientWithResponsesInterface(ctrl)

			m.EXPECT().GetProviderWithResponse(gomock.Any(), gomock.Eq("commonfate"), gomock.Eq("test"), gomock.Eq(tc.version)).Return(&tc.mockGetProviderRes, tc.mockGetProviderResErr)

			s := Service{
				Clock:                  clk,
				DB:                     dbc,
				ProviderRegistryClient: m,
			}

			got, err := s.CreateTargetGroup(context.Background(), tc.give)

			if err != nil {
				assert.Equal(t, tc.wantErr, err)
			}
			assert.Equal(t, tc.want, got)

		})
	}
}
