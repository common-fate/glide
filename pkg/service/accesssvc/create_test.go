package accesssvc

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/benbjohnson/clock"
	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/ddb"
	"github.com/common-fate/ddb/ddbmock"
	ahTypes "github.com/common-fate/granted-approvals/accesshandler/pkg/types"
	"github.com/common-fate/granted-approvals/pkg/access"
	"github.com/common-fate/granted-approvals/pkg/identity"
	"github.com/common-fate/granted-approvals/pkg/rule"
	accessMocks "github.com/common-fate/granted-approvals/pkg/service/accesssvc/mocks"
	"github.com/common-fate/granted-approvals/pkg/storage"
	"github.com/common-fate/granted-approvals/pkg/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewRequest(t *testing.T) {
	type createGrantResponse struct {
		request *access.Request
		err     error
	}
	type testcase struct {
		name                    string
		giveInput               types.CreateRequestRequest
		giveUser                identity.User
		rule                    *rule.AccessRule
		ruleErr                 error
		wantErr                 error
		want                    *CreateRequestResult
		withCreateGrantResponse createGrantResponse
		withGetGroupResponse    *storage.GetGroup
		validationResponse      []ahTypes.GrantValidation
	}

	clk := clock.NewMock()
	autoApproval := types.AUTOMATIC
	reviewed := types.REVIEWED
	testcases := []testcase{
		{
			name: "ok, no approvers so should auto approve",
			//just passing the group here, technically a user isnt an approver
			giveUser: identity.User{Groups: []string{"a"}},
			rule: &rule.AccessRule{
				Groups: []string{"a"},
			},
			want: &CreateRequestResult{
				Request: access.Request{
					ID:             "-",
					Status:         access.APPROVED,
					CreatedAt:      clk.Now(),
					UpdatedAt:      clk.Now(),
					Grant:          &access.Grant{},
					ApprovalMethod: &autoApproval,
					SelectedWith:   make(map[string]access.Option),
				},
			},
			withCreateGrantResponse: createGrantResponse{
				request: &access.Request{
					ID:             "-",
					Status:         access.APPROVED,
					CreatedAt:      clk.Now(),
					UpdatedAt:      clk.Now(),
					Grant:          &access.Grant{},
					ApprovalMethod: &autoApproval,
					SelectedWith:   make(map[string]access.Option),
				},
			},
			validationResponse: []ahTypes.GrantValidation{},
		},
		{
			name:     "fails because requested duration is greater than max duration",
			giveUser: identity.User{Groups: []string{"a"}},
			giveInput: types.CreateRequestRequest{
				Timing: types.RequestTiming{
					DurationSeconds: 20,
				},
			},
			rule: &rule.AccessRule{
				Groups: []string{"a"},
				TimeConstraints: types.TimeConstraints{
					MaxDurationSeconds: 10,
				},
			},
			wantErr: &apio.APIError{
				Err:    errors.New("request validation failed"),
				Status: http.StatusBadRequest,
				Fields: []apio.FieldError{
					{
						Field: "timing.durationSeconds",
						Error: fmt.Sprintf("durationSeconds: %d exceeds the maximum duration seconds: %d", 20, 10),
					},
				},
			},
			validationResponse: []ahTypes.GrantValidation{},
		},
		{
			name:     "user not in correct group",
			giveUser: identity.User{Groups: []string{"a"}},
			rule: &rule.AccessRule{
				Groups: []string{"b"},
			},
			wantErr: ErrNoMatchingGroup,
		},
		{
			name:               "rule not found",
			giveUser:           identity.User{Groups: []string{"a"}},
			ruleErr:            ddb.ErrNoItems,
			wantErr:            ErrRuleNotFound,
			validationResponse: []ahTypes.GrantValidation{},
		},
		{
			name:     "with reviewers",
			giveUser: identity.User{Groups: []string{"a"}},
			rule: &rule.AccessRule{
				Groups: []string{"a"},
				Approval: rule.Approval{
					Users: []string{"b"},
				},
			},
			want: &CreateRequestResult{
				Request: access.Request{
					ID:             "-",
					Status:         access.PENDING,
					CreatedAt:      clk.Now(),
					UpdatedAt:      clk.Now(),
					ApprovalMethod: &reviewed,
					SelectedWith:   make(map[string]access.Option),
				},
				Reviewers: []access.Reviewer{
					{
						ReviewerID: "b",
						Request: access.Request{
							ID:             "-",
							Status:         access.PENDING,
							CreatedAt:      clk.Now(),
							UpdatedAt:      clk.Now(),
							ApprovalMethod: &reviewed,
							SelectedWith:   make(map[string]access.Option),
						},
					},
				},
			},
		},
		{
			name:     "requestor is approver on access rule",
			giveUser: identity.User{ID: "a", Groups: []string{"a"}},
			rule: &rule.AccessRule{
				Groups: []string{"a"},
				Approval: rule.Approval{
					Users: []string{"a", "b"},
				},
			},
			// user 'a' should not be included as an approver of this request,
			// as they made the request.
			want: &CreateRequestResult{
				Request: access.Request{
					ID:             "-",
					RequestedBy:    "a",
					Status:         access.PENDING,
					CreatedAt:      clk.Now(),
					UpdatedAt:      clk.Now(),
					ApprovalMethod: &reviewed,
					SelectedWith:   make(map[string]access.Option),
				},
				Reviewers: []access.Reviewer{
					{
						ReviewerID: "b",
						Request: access.Request{
							ID:             "-",
							RequestedBy:    "a",
							Status:         access.PENDING,
							CreatedAt:      clk.Now(),
							UpdatedAt:      clk.Now(),
							ApprovalMethod: &reviewed,
							SelectedWith:   make(map[string]access.Option),
						},
					},
				},
			},
			validationResponse: []ahTypes.GrantValidation{},
		},
		{
			name:     "requestor is in approver group on access rule",
			giveUser: identity.User{ID: "a", Groups: []string{"a", "b"}},
			rule: &rule.AccessRule{
				Groups: []string{"a"},
				Approval: rule.Approval{
					Groups: []string{"b"},
				},
			},
			withGetGroupResponse: &storage.GetGroup{
				Result: &identity.Group{
					ID:    "b",
					Users: []string{"c"},
				},
			},
			// user 'a' should not be included as an approver of this request,
			// as they made the request.
			want: &CreateRequestResult{
				Request: access.Request{
					ID:             "-",
					RequestedBy:    "a",
					Status:         access.PENDING,
					CreatedAt:      clk.Now(),
					UpdatedAt:      clk.Now(),
					ApprovalMethod: &reviewed,
					SelectedWith:   make(map[string]access.Option),
				},
				Reviewers: []access.Reviewer{
					{
						ReviewerID: "c",
						Request: access.Request{
							ID:             "-",
							RequestedBy:    "a",
							Status:         access.PENDING,
							CreatedAt:      clk.Now(),
							UpdatedAt:      clk.Now(),
							ApprovalMethod: &reviewed,
							SelectedWith:   make(map[string]access.Option),
						},
					},
				},
			},
			validationResponse: []ahTypes.GrantValidation{},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			db := ddbmock.New(t)
			db.MockQueryWithErr(&storage.GetAccessRuleCurrent{Result: tc.rule}, tc.ruleErr)
			db.MockQuery(tc.withGetGroupResponse)
			db.MockQuery(&storage.ListRequestReviewers{})
			db.MockQuery(&storage.ListRequestsForUserAndRuleAndRequestend{})
			ctrl := gomock.NewController(t)

			defer ctrl.Finish()

			ctrl2 := gomock.NewController(t)
			ep := accessMocks.NewMockEventPutter(ctrl2)
			ep.EXPECT().Put(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

			g := accessMocks.NewMockGranter(ctrl)
			g.EXPECT().CreateGrant(gomock.Any(), gomock.Any()).Return(tc.withCreateGrantResponse.request, tc.withCreateGrantResponse.err).AnyTimes()
			g.EXPECT().ValidateGrant(gomock.Any(), gomock.Any()).Return(tc.validationResponse, nil).AnyTimes()

			ca := accessMocks.NewMockCacheService(ctrl)
			ca.EXPECT().LoadCachedProviderArgOptions(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, nil, nil).AnyTimes()
			s := Service{
				Clock:       clk,
				DB:          db,
				Granter:     g,
				EventPutter: ep,
				Cache:       ca,
			}
			got, err := s.CreateRequest(context.Background(), &tc.giveUser, tc.giveInput)
			if got != nil {
				// ignore the autogenerated ID for testing.
				got.Request.ID = "-"

				for i := range got.Reviewers {
					got.Reviewers[i].Request.ID = "-"
				}
			}

			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.want, got)
		})
	}

}
