package rulesvc

import (
	"testing"

	"context"

	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/ddb/ddbmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetRule(t *testing.T) {
	type testcase struct {
		name            string
		givenUser       identity.User
		getRuleResponse *rule.AccessRule
		wantErr         error
		want            *rule.AccessRuleWithCanRequest
		isAdmin         bool
	}

	var mockReq = &rule.AccessRule{
		ID:     "rule1",
		Groups: []string{"group1"},
		Approval: rule.Approval{
			Users: []string{"approver1"},
		},
	}

	testcases := []testcase{
		{
			name:            "User can see a rule they're an approver for",
			givenUser:       identity.User{ID: "approver1", Groups: []string{"group1"}},
			getRuleResponse: mockReq,
			want: &rule.AccessRuleWithCanRequest{
				Rule:       mockReq,
				CanRequest: true,
			},
		},
		{
			name:            "Approval cannot request a rule if they are not in request group",
			givenUser:       identity.User{ID: "approver1", Groups: []string{}},
			getRuleResponse: mockReq,
			want: &rule.AccessRuleWithCanRequest{
				Rule:       mockReq,
				CanRequest: false,
			},
		},
		{
			name:            "User can see a rule they're assigned to (via the groups)",
			givenUser:       identity.User{ID: "approver2", Groups: []string{"group1"}},
			getRuleResponse: mockReq,
			want: &rule.AccessRuleWithCanRequest{
				Rule:       mockReq,
				CanRequest: true,
			},
		},
		{
			name:            "User *cannot* see if they're neither an approver nor assigned to the rule",
			givenUser:       identity.User{ID: "a", Groups: []string{"group1"}},
			getRuleResponse: mockReq,
			wantErr:         ErrUserNotAuthorized,
		},
		{
			name:            "Admins can always access rules",
			givenUser:       identity.User{ID: "a"},
			isAdmin:         true,
			getRuleResponse: mockReq,
			want: &rule.AccessRuleWithCanRequest{
				Rule:       mockReq,
				CanRequest: false,
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			db := ddbmock.New(t)
			db.MockQueryWithErr(&storage.GetAccessRuleCurrent{Result: tc.getRuleResponse}, tc.wantErr)

			ctrl := gomock.NewController(t)

			defer ctrl.Finish()

			s := Service{
				DB: db,
			}
			got, err := s.GetRule(context.Background(), tc.getRuleResponse.ID, &tc.givenUser, tc.isAdmin)

			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.want, got)
		})
	}

}
