package workflowsvc

// import (
// 	"context"
// 	"testing"
// 	"time"

// 	"github.com/benbjohnson/clock"
// 	"github.com/common-fate/common-fate/pkg/identity"
// 	"github.com/common-fate/common-fate/pkg/requests"
// 	"github.com/common-fate/common-fate/pkg/rule"
// 	"github.com/common-fate/common-fate/pkg/service/workflowsvc/mocks"
// 	"github.com/common-fate/common-fate/pkg/storage"
// 	"github.com/common-fate/common-fate/pkg/types"
// 	"github.com/common-fate/ddb"
// 	"github.com/common-fate/ddb/ddbmock"
// 	"github.com/golang/mock/gomock"
// 	"github.com/stretchr/testify/assert"
// )

// func TestRevokeGrant(t *testing.T) {
// 	type testcase struct {
// 		name                       string
// 		withRevokeGrantResponseErr error
// 		Grants                     []requests.Grantv2
// 		AccessGroups               []requests.AccessGroup
// 		getRule                    rule.AccessRule
// 		withGetRuleErr             error
// 		giveRequest                requests.Requestv2
// 		wantGrantErr               error
// 		wantGroupErr               error
// 		want                       *requests.Requestv2
// 		revokerID                  string
// 		// requestReviewers              []access.Reviewer
// 		subject string
// 	}
// 	clk := clock.NewMock()

// 	testcases := []testcase{
// 		{
// 			name: "ok",
// 			getRule: rule.AccessRule{ID: "rule1",
// 				Status: rule.ACTIVE,

// 				Description: "string",
// 				Name:        "string",
// 				Groups:      []string{"string"},
// 				// Target: rule.Target{
// 				// 	With:          map[string]string{},
// 				// 	TargetGroupID: "123",
// 				// },
// 			},
// 			giveRequest: requests.Requestv2{
// 				RequestedBy: identity.User{ID: "123"},
// 			},
// 			withRevokeGrantResponseErr: nil,
// 			Grants: []requests.Grantv2{
// 				{
// 					ID:          "gra_123",
// 					AccessGroup: "123",
// 					Status:      types.GrantStatusACTIVE,
// 					End:         clk.Now().Add(time.Hour * 2),
// 				},
// 			},
// 			AccessGroups: []requests.AccessGroup{
// 				{ID: "123"},
// 			},
// 			want: &requests.Requestv2{
// 				RequestedBy: identity.User{ID: "123"},
// 			},
// 			revokerID: "test",
// 			// requestReviewers:           []access.Reviewer{{ReviewerID: "123"}},
// 			subject:        "test@commonfate.io",
// 			withGetRuleErr: nil,
// 		},
// 		{
// 			name: "no grant",
// 			getRule: rule.AccessRule{ID: "rule1",
// 				Status: rule.ACTIVE,

// 				Description: "string",
// 				Name:        "string",
// 				Groups:      []string{"string"},
// 				// Target: rule.Target{
// 				// 	With:          map[string]string{},
// 				// 	TargetGroupID: "123",
// 				// },
// 			},
// 			giveRequest: requests.Requestv2{
// 				RequestedBy: identity.User{ID: "123"},
// 			},
// 			withRevokeGrantResponseErr: ErrNoGrant,
// 			Grants:                     nil,
// 			AccessGroups: []requests.AccessGroup{
// 				{ID: "123"},
// 			},
// 			want:         nil,
// 			wantGrantErr: ddb.ErrNoItems,
// 			wantGroupErr: nil,
// 			revokerID:    "test",
// 			// requestReviewers:           []access.Reviewer{{ReviewerID: "123"}},
// 			subject:        "test@commonfate.io",
// 			withGetRuleErr: nil,
// 		},
// 		{
// 			name: "no access group",
// 			getRule: rule.AccessRule{ID: "rule1",
// 				Status: rule.ACTIVE,

// 				Description: "string",
// 				Name:        "string",
// 				Groups:      []string{"string"},
// 				// Target: rule.Target{
// 				// 	With:          map[string]string{},
// 				// 	TargetGroupID: "123",
// 				// },
// 			},
// 			giveRequest: requests.Requestv2{
// 				RequestedBy: identity.User{ID: "123"},
// 			},
// 			withRevokeGrantResponseErr: ddb.ErrNoItems,
// 			Grants:                     nil,
// 			AccessGroups: []requests.AccessGroup{
// 				{ID: "123"},
// 			},
// 			want:         nil,
// 			wantGroupErr: ddb.ErrNoItems,
// 			revokerID:    "test",
// 			// requestReviewers:           []access.Reviewer{{ReviewerID: "123"}},
// 			subject:        "test@commonfate.io",
// 			withGetRuleErr: nil,
// 		},

// 		{
// 			name: "trying to revoke inactive grant",
// 			getRule: rule.AccessRule{ID: "rule1",
// 				Status: rule.ACTIVE,

// 				Description: "string",
// 				Name:        "string",
// 				Groups:      []string{"string"},
// 				// Target: rule.Target{
// 				// 	With:          map[string]string{},
// 				// 	TargetGroupID: "123",
// 				// },
// 			},
// 			giveRequest: requests.Requestv2{
// 				RequestedBy: identity.User{ID: "123"},
// 			},
// 			withRevokeGrantResponseErr: ErrGrantInactive,
// 			Grants: []requests.Grantv2{
// 				{
// 					ID:          "gra_123",
// 					AccessGroup: "123",
// 					Status:      types.GrantStatusEXPIRED,
// 					End:         clk.Now().Add(time.Hour * 2),
// 				},
// 			},
// 			AccessGroups: []requests.AccessGroup{
// 				{ID: "123"},
// 			},
// 			want:      nil,
// 			revokerID: "test",
// 			// requestReviewers:           []access.Reviewer{{ReviewerID: "123"}},
// 			subject:        "test@commonfate.io",
// 			withGetRuleErr: nil,
// 		},
// 	}

// 	for _, tc := range testcases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			runtime := mocks.NewMockRuntime(ctrl)
// 			runtime.EXPECT().Revoke(gomock.Any(), gomock.Any()).Return(tc.withRevokeGrantResponseErr).AnyTimes()

// 			eventbus := mocks.NewMockEventPutter(ctrl)
// 			eventbus.EXPECT().Put(gomock.Any(), gomock.Any()).Return(tc.withRevokeGrantResponseErr).AnyTimes()

// 			c := ddbmock.New(t)
// 			c.MockQueryWithErr(&storage.ListAccessGroups{Result: tc.AccessGroups}, tc.wantGroupErr)
// 			c.MockQueryWithErr(&storage.ListGrantsV2{Result: tc.Grants}, tc.wantGrantErr)
// 			c.MockQueryWithErr(&storage.GetAccessRule{Result: &tc.getRule}, tc.withGetRuleErr)

// 			s := Service{
// 				Runtime:  runtime,
// 				DB:       c,
// 				Clk:      clk,
// 				Eventbus: eventbus,
// 			}

// 			gotRequest, err := s.Revoke(context.Background(), tc.giveRequest, tc.revokerID, tc.subject)
// 			assert.Equal(t, tc.withRevokeGrantResponseErr, err)

// 			assert.Equal(t, tc.want, gotRequest)
// 		})
// 	}
// }
