package workflowsvc

import (
	"context"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	ahTypes "github.com/common-fate/common-fate/accesshandler/pkg/types"
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/common-fate/common-fate/pkg/service/workflowsvc/mocks"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/ddb"
	"github.com/common-fate/ddb/ddbmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRevokeGrant(t *testing.T) {
	type testcase struct {
		name                          string
		withRevokeGrantResponseErr    error
		withUser                      *identity.User
		getRule                       rule.AccessRule
		giveRequest                   access.Request
		wantUserErr                   error
		want                          *access.Request
		revokerID                     string
		requestReviewers              []access.Reviewer
		withGetRuleVersionResponseErr error
	}
	clk := clock.NewMock()

	testcases := []testcase{
		{
			name: "ok",
			getRule: rule.AccessRule{ID: "rule1",
				Status: rule.ACTIVE,

				Description: "string",
				Name:        "string",
				Groups:      []string{"string"},
				Target: rule.Target{
					ProviderID:    "string",
					With:          map[string]string{},
					TargetGroupID: "123",
				}},
			giveRequest: access.Request{
				RequestedBy: "user1",
				Grant: &access.Grant{
					Status: ahTypes.GrantStatus(ahTypes.GrantStatusACTIVE),
					End:    time.Now().Add(time.Hour),
				},
			},
			withRevokeGrantResponseErr: nil,
			withUser:                   &identity.User{Groups: []string{"testAdmin"}},
			wantUserErr:                nil,
			want:                       nil,
			revokerID:                  "test",
			requestReviewers:           []access.Reviewer{{ReviewerID: "123"}},
		},
		{
			name: "no grant",
			getRule: rule.AccessRule{ID: "rule1",
				Status: rule.ACTIVE,

				Description: "string",
				Name:        "string",
				Groups:      []string{"string"},
				Target: rule.Target{
					ProviderID:    "string",
					With:          map[string]string{},
					TargetGroupID: "123",
				}},
			giveRequest: access.Request{
				RequestedBy: "user1",
				Grant: &access.Grant{
					Status: ahTypes.GrantStatus(ahTypes.GrantStatusACTIVE),
					End:    time.Now().Add(time.Hour),
				},
			},
			withRevokeGrantResponseErr: ErrNoGrant,
			withUser:                   &identity.User{Groups: []string{"testAdmin"}},
			wantUserErr:                nil,
			want:                       nil,
			revokerID:                  "test",
			requestReviewers:           []access.Reviewer{{ReviewerID: "123"}},
		},
		{
			name: "trying to revoke inactive grant",
			getRule: rule.AccessRule{ID: "rule1",
				Status: rule.ACTIVE,

				Description: "string",
				Name:        "string",
				Groups:      []string{"string"},
				Target: rule.Target{
					ProviderID:    "string",
					With:          map[string]string{},
					TargetGroupID: "123",
				}},
			giveRequest: access.Request{
				RequestedBy: "user1",
				Grant: &access.Grant{
					Status: ahTypes.GrantStatusEXPIRED,
				},
			},
			withRevokeGrantResponseErr: ErrGrantInactive,
			withUser:                   &identity.User{Groups: []string{"testAdmin"}},
			wantUserErr:                nil,
			want:                       nil,
			revokerID:                  "test",
		},
		{
			name: "access rule version not found",
			getRule: rule.AccessRule{ID: "rule1",
				Status: rule.ACTIVE,

				Description: "string",
				Name:        "string",
				Groups:      []string{"string"},
				Target: rule.Target{
					ProviderID:    "string",
					With:          map[string]string{},
					TargetGroupID: "123",
				}},
			giveRequest: access.Request{
				RequestedBy: "user1",
				Grant: &access.Grant{
					Status: ahTypes.GrantStatus(ahTypes.GrantStatusACTIVE),
					End:    time.Now().Add(time.Hour),
				},
			},
			withRevokeGrantResponseErr:    ddb.ErrNoItems,
			withUser:                      &identity.User{Groups: []string{"testAdmin"}},
			wantUserErr:                   nil,
			want:                          nil,
			revokerID:                     "test",
			requestReviewers:              []access.Reviewer{{ReviewerID: "123"}},
			withGetRuleVersionResponseErr: ddb.ErrNoItems,
		},
	}

	for i := range testcases {
		tc := testcases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			runtime := mocks.NewMockRuntime(ctrl)
			runtime.EXPECT().Revoke(gomock.Any(), gomock.Any(), gomock.Any()).Return(tc.withRevokeGrantResponseErr).AnyTimes()

			eventbus := mocks.NewMockEventPutter(ctrl)
			eventbus.EXPECT().Put(gomock.Any(), gomock.Any()).Return(tc.withRevokeGrantResponseErr).AnyTimes()

			c := ddbmock.New(t)
			c.MockQueryWithErr(&storage.GetUser{Result: tc.withUser}, tc.wantUserErr)
			c.MockQueryWithErr(&storage.GetAccessRuleVersion{Result: &tc.getRule}, tc.withGetRuleVersionResponseErr)
			c.MockQueryWithErr(&storage.ListRequestReviewers{Result: tc.requestReviewers}, tc.wantUserErr)

			s := Service{
				Runtime:  runtime,
				DB:       c,
				Clk:      clk,
				Eventbus: eventbus,
			}

			gotRequest, err := s.Revoke(context.Background(), tc.giveRequest, tc.revokerID, tc.withUser.Email)
			assert.Equal(t, tc.withRevokeGrantResponseErr, err)

			assert.Equal(t, tc.want, gotRequest)
		})
	}
}
