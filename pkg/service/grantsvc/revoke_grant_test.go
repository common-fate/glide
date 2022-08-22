package grantsvc

import (
	"context"
	"testing"
	"time"

	"github.com/benbjohnson/clock"

	ah_types "github.com/common-fate/granted-approvals/accesshandler/pkg/types"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/types/ahmocks"

	"github.com/common-fate/granted-approvals/pkg/access"
	"github.com/common-fate/iso8601"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAccessRevoke(t *testing.T) {

	type testcase struct {
		name                    string
		give                    RevokeGrantOpts
		wantErr                 error
		withRevokeGrantResponse ah_types.PostGrantsRevokeResponse
	}
	clk := clock.NewMock()

	testStartTime := iso8601.Now()
	testEndTime := iso8601.Now().Add(time.Hour)

	testcases := []testcase{

		{
			name: "Trying to revoke inactive grant",

			withRevokeGrantResponse: ah_types.PostGrantsRevokeResponse{
				JSON200: &struct {
					AdditionalProperties ah_types.AdditionalProperties "json:\"additionalProperties\""
					Grant                ah_types.Grant                "json:\"grant\""
				}{AdditionalProperties: ah_types.AdditionalProperties{},
					Grant: ah_types.Grant{
						ID:      "123",
						Start:   iso8601.New(testStartTime.Time),
						End:     iso8601.New(testEndTime.Add(time.Minute * 2)),
						Subject: "test@test.com",
						Status:  "REVOKED",
					}},
			},
			wantErr: ErrGrantInactive,

			give: RevokeGrantOpts{Request: access.Request{
				ID: "123",
				Grant: &access.Grant{
					Start:    testStartTime.Time,
					End:      testEndTime,
					Subject:  "test@test.com",
					Status:   "PENDING",
					Provider: "okta",
				}}, RevokerID: "1234"}},
	}

	for _, tc := range testcases {

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			g := ahmocks.NewMockClientWithResponsesInterface(ctrl)
			g.EXPECT().PostGrantsRevokeWithResponse(gomock.Any(), "123", ah_types.PostGrantsRevokeJSONRequestBody{
				RevokerId: tc.give.RevokerID,
			}).Return(&tc.withRevokeGrantResponse, tc.wantErr).AnyTimes()

			s := Granter{AHClient: g, Clock: clk}
			_, err := s.RevokeGrant(context.Background(), tc.give)

			assert.Equal(t, tc.wantErr, err)
			//assert.Equal(t, tc.wantResp, gotGrant)
		})
	}

}
