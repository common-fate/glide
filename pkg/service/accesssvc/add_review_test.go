package accesssvc

import (
	"context"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/storage"

	"github.com/common-fate/common-fate/pkg/service/accesssvc/mocks"
	"github.com/common-fate/common-fate/pkg/service/grantsvc"
	"github.com/common-fate/ddb/ddbmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

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
			c.MockQuery(&storage.ListRequestsForUserAndRequestend{})

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
