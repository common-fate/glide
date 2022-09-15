package accesssvc

import (
	"context"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/common-fate/granted-approvals/pkg/access"
	"github.com/common-fate/granted-approvals/pkg/storage"

	"github.com/common-fate/ddb/ddbmock"
	"github.com/common-fate/granted-approvals/pkg/service/accesssvc/mocks"
	"github.com/common-fate/granted-approvals/pkg/service/grantsvc"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestOverlapsExistingGrant(t *testing.T) {
	type testcase struct {
		name           string
		grant          access.Grant
		existingGrants []access.Request
		want           bool
	}
	clk := clock.NewMock()
	a := access.Grant{Start: clk.Now().Add(-time.Minute), End: clk.Now().Add(-time.Minute + time.Second)}
	b := access.Grant{Start: clk.Now(), End: clk.Now().Add(time.Minute)}
	c := access.Grant{Start: clk.Now().Add(-time.Minute), End: clk.Now().Add(time.Minute)}
	d := access.Grant{Start: clk.Now().Add(time.Second * 30), End: clk.Now().Add(time.Minute)}

	testcases := []testcase{
		{
			name:           "no existing grants",
			grant:          a,
			existingGrants: []access.Request{},
			want:           false,
		},
		{
			name:           "no overlap before",
			grant:          a,
			existingGrants: []access.Request{{Grant: &b}},
			want:           false,
		},
		{
			name:           "overlap",
			grant:          a,
			existingGrants: []access.Request{{Grant: &a}},
			want:           true,
		},
		{
			name:           "partial overlap",
			grant:          c,
			existingGrants: []access.Request{{Grant: &b}},
			want:           true,
		},
		{
			name:           "partial overlap",
			grant:          d,
			existingGrants: []access.Request{{Grant: &b}},
			want:           true,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := overlapsExistingGrant(tc.grant.Start, tc.grant.End, tc.existingGrants)
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
			g.EXPECT().ValidateGrant(gomock.Any(), gomock.Any()).Return(tc.withCreateGrantResponse.err).AnyTimes()

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
