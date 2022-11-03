package accesssvc

import (
	"context"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/common-fate/granted-approvals/pkg/access"
	"github.com/common-fate/granted-approvals/pkg/rule"
	"github.com/common-fate/granted-approvals/pkg/storage"

	"github.com/common-fate/ddb/ddbmock"
	"github.com/common-fate/granted-approvals/pkg/service/accesssvc/mocks"
	"github.com/common-fate/granted-approvals/pkg/service/grantsvc"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestOverlapsExistingGrant(t *testing.T) {
	type testcase struct {
		name               string
		accessRequest      access.Request
		upcomingRequests   []access.Request
		currentRequestRule rule.AccessRule
		allRules           []rule.AccessRule
		want               bool
		clock              clock.Clock
	}
	clk := clock.NewMock()
	inOneMinute := clk.Now().Add(time.Minute)

	now := clk.Now()

	//some requests premade with specific timings
	a := access.Request{RequestedTiming: access.Timing{StartTime: &inOneMinute, Duration: time.Minute * 2}, Rule: "rule_a"} //started 1 minute ago ends in a minute
	b := access.Request{RequestedTiming: access.Timing{StartTime: &now, Duration: time.Minute * 2}, Rule: "rule_a"}         //started now, ends in 2 minute
	args1 := make(map[string]string)
	args1["1"] = "arg1"
	args1["2"] = "arg2"

	args2 := make(map[string]string)
	args1["a"] = "argA"
	args1["b"] = "argB"

	testcases := []testcase{
		{
			name:               "no existing grants",
			accessRequest:      access.Request{ID: "123"},
			upcomingRequests:   []access.Request{},
			currentRequestRule: rule.AccessRule{Target: rule.Target{}},
			allRules:           []rule.AccessRule{},
			clock:              clk,
			want:               false,
		},
		{
			name:               "different provider passes",
			accessRequest:      access.Request{ID: "abc", Rule: "rule_b"},
			upcomingRequests:   []access.Request{{ID: "def", Rule: "rule_a"}},
			currentRequestRule: rule.AccessRule{ID: "rule_a", Target: rule.Target{ProviderID: "prov_a"}},
			allRules:           []rule.AccessRule{{ID: "rule_a", Target: rule.Target{ProviderID: "prov_b"}}},
			clock:              clk,
			want:               false,
		},
		{
			name:               "request overlaps current active request fails",
			accessRequest:      access.Request{ID: "123", Rule: "rule_a", RequestedTiming: access.Timing{StartTime: &now, Duration: time.Minute * 5}},
			upcomingRequests:   []access.Request{a},
			currentRequestRule: rule.AccessRule{ID: "rule_a", Target: rule.Target{ProviderID: "prov_a"}},
			allRules:           []rule.AccessRule{{ID: "rule_a", Target: rule.Target{ProviderID: "prov_a"}}},
			clock:              clk,
			want:               true,
		},
		{
			name:               "scheduled request overlaps current active request fails",
			accessRequest:      access.Request{ID: "123", Rule: "rule_a", RequestedTiming: access.Timing{StartTime: &inOneMinute, Duration: time.Minute * 5}},
			upcomingRequests:   []access.Request{b},
			currentRequestRule: rule.AccessRule{ID: "rule_a", Target: rule.Target{ProviderID: "prov_a"}},
			allRules:           []rule.AccessRule{{ID: "rule_a", Target: rule.Target{ProviderID: "prov_a"}}},
			clock:              clk,
			want:               true,
		},
		{
			name:               "same rule different arguments should succeed",
			accessRequest:      access.Request{ID: "123", Rule: "rule_a", RequestedTiming: access.Timing{StartTime: &inOneMinute, Duration: time.Minute * 5}},
			upcomingRequests:   []access.Request{b},
			currentRequestRule: rule.AccessRule{ID: "rule_a", Target: rule.Target{ProviderID: "prov_a", With: args1}},
			allRules:           []rule.AccessRule{{ID: "rule_a", Target: rule.Target{ProviderID: "prov_a", With: args2}}},
			clock:              clk,
			want:               false,
		},
		{
			name:               "same rule same arguments should fail",
			accessRequest:      access.Request{ID: "123", Rule: "rule_a", RequestedTiming: access.Timing{StartTime: &inOneMinute, Duration: time.Minute * 5}},
			upcomingRequests:   []access.Request{b},
			currentRequestRule: rule.AccessRule{ID: "rule_a", Target: rule.Target{ProviderID: "prov_a", With: args1}},
			allRules:           []rule.AccessRule{{ID: "rule_a", Target: rule.Target{ProviderID: "prov_a", With: args1}}},
			clock:              clk,
			want:               true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := overlapsExistingGrantCheck(tc.accessRequest, tc.upcomingRequests, tc.currentRequestRule, tc.allRules, tc.clock)
			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})

	}

}

func TestAddReview(t *testing.T) {
	type createGrantResponse struct {
		request *access.Request
		err     error
	}
	type testcase struct {
		name                    string
		give                    AddReviewOpts
		want                    *AddReviewResult
		wantErr                 error
		withCreateGrantResponse createGrantResponse
		wantCreateGrantOpts     grantsvc.CreateGrantOpts
	}

	clk := clock.NewMock()

	now := clk.Now()
	overrideTiming := &access.Timing{
		Duration:  time.Minute,
		StartTime: &now,
	}
	requestWithOverride := access.Request{
		Status:         access.APPROVED,
		Grant:          &access.Grant{},
		OverrideTiming: overrideTiming,
		UpdatedAt:      clk.Now(),
	}
	testcases := []testcase{
		{
			name: "ok",
			give: AddReviewOpts{
				ReviewerID: "a",
				Decision:   access.DecisionApproved,
				Reviewers: []access.Reviewer{
					{
						ReviewerID: "a",
						Request: access.Request{
							Status: access.PENDING,
						},
					},
					{
						ReviewerID: "b",
						Request: access.Request{
							Status: access.PENDING,
						},
					},
				},
				Request: access.Request{
					Status: access.PENDING,
				},
			},
			wantCreateGrantOpts: grantsvc.CreateGrantOpts{
				Request: access.Request{
					Status: access.APPROVED,
				},
			},
			withCreateGrantResponse: createGrantResponse{

				request: &access.Request{
					Status:    access.APPROVED, // request should be approved
					UpdatedAt: clk.Now(),
					Grant:     &access.Grant{},
				},
			},
			want: &AddReviewResult{
				Request: access.Request{
					Status:    access.APPROVED, // request should be approved
					UpdatedAt: clk.Now(),
					Grant:     &access.Grant{},
				},
			},
		},
		{
			name: "ok with override times",
			give: AddReviewOpts{
				ReviewerID: "a",
				Decision:   access.DecisionApproved,
				Reviewers: []access.Reviewer{
					{
						ReviewerID: "a",
						Request: access.Request{
							Status: access.PENDING,
						},
					},
					{
						ReviewerID: "b",
						Request: access.Request{
							Status: access.PENDING,
						},
					},
				},
				Request: access.Request{
					Status: access.PENDING,
				},
				OverrideTiming: overrideTiming,
			},
			wantCreateGrantOpts: grantsvc.CreateGrantOpts{
				Request: access.Request{
					Status:         access.APPROVED,
					OverrideTiming: overrideTiming,
				},
			},
			withCreateGrantResponse: createGrantResponse{

				request: &requestWithOverride,
			},
			want: &AddReviewResult{
				Request: requestWithOverride,
			},
		},
		{
			name: "cannot review own request",
			give: AddReviewOpts{
				ReviewerID: "a",
				Decision:   access.DecisionApproved,
				Reviewers: []access.Reviewer{

					{
						ReviewerID: "b",
						Request: access.Request{
							Status: access.PENDING,
						},
					},
				},
				Request: access.Request{
					Status:      access.PENDING,
					RequestedBy: "a",
				},
			},
			wantErr: ErrUserNotAuthorized,
		},
		{
			name: "admin cannot review own request",
			give: AddReviewOpts{
				ReviewerID:      "a",
				Decision:        access.DecisionApproved,
				ReviewerIsAdmin: true,
				Request: access.Request{
					Status:      access.PENDING,
					RequestedBy: "a",
				},
			},
			wantErr: ErrUserNotAuthorized,
		},
		{
			name: "admin can review not own request",
			give: AddReviewOpts{
				ReviewerID:      "a",
				Decision:        access.DecisionApproved,
				ReviewerIsAdmin: true,
				Request: access.Request{
					Status:      access.PENDING,
					RequestedBy: "b",
				},
			},
			wantCreateGrantOpts: grantsvc.CreateGrantOpts{
				Request: access.Request{
					Status:      access.APPROVED,
					RequestedBy: "b",
				},
			},
			withCreateGrantResponse: createGrantResponse{

				request: &access.Request{
					Status:      access.APPROVED, // request should be approved
					RequestedBy: "b",
					UpdatedAt:   clk.Now(),
					Grant:       &access.Grant{},
				},
			},
			want: &AddReviewResult{
				Request: access.Request{
					Status:      access.APPROVED, // request should be approved
					RequestedBy: "b",
					UpdatedAt:   clk.Now(),
					Grant:       &access.Grant{},
				},
			},
		},
	}

	for _, tc := range testcases {

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			g := mocks.NewMockGranter(ctrl)
			g.EXPECT().CreateGrant(gomock.Any(), gomock.Eq(tc.wantCreateGrantOpts)).Return(tc.withCreateGrantResponse.request, tc.withCreateGrantResponse.err).AnyTimes()

			ctrl2 := gomock.NewController(t)
			ep := mocks.NewMockEventPutter(ctrl2)
			ep.EXPECT().Put(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

			c := ddbmock.New(t)
			c.MockQuery(&storage.ListRequestsForUserAndRuleAndRequestend{})

			// called by dbupdate.GetUpdateRequestItems
			c.MockQuery(&storage.ListRequestReviewers{})

			s := Service{
				Clock:       clk,
				DB:          c,
				Granter:     g,
				EventPutter: ep,
			}
			got, err := s.AddReviewAndGrantAccess(context.Background(), tc.give)
			if tc.wantErr == nil {
				assert.NoError(t, err)
			}
			if tc.wantErr != nil {
				assert.EqualError(t, err, tc.wantErr.Error())
			}
			assert.Equal(t, tc.want, got)
		})
	}

}
