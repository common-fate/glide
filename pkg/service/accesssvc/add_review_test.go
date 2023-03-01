package accesssvc

import (
	"context"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/service/accesssvc/mocks"
	accessMocks "github.com/common-fate/common-fate/pkg/service/accesssvc/mocks"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/types"
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
	}

	clk := clock.NewMock()

	now := clk.Now()
	overrideTiming := &access.Timing{
		Duration:  time.Minute,
		StartTime: &now,
	}
	reviewed := types.REVIEWED
	requestWithOverride := access.Request{
		Status:         access.APPROVED,
		Grant:          &access.Grant{},
		OverrideTiming: overrideTiming,
		UpdatedAt:      clk.Now(),
		ApprovalMethod: &reviewed,
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
			withCreateGrantResponse: createGrantResponse{

				request: &access.Request{
					Status:    access.APPROVED, // request should be approved
					UpdatedAt: clk.Now(),
					Grant:     &access.Grant{},
				},
			},
			want: &AddReviewResult{
				Request: access.Request{
					Status:         access.APPROVED, // request should be approved
					UpdatedAt:      clk.Now(),
					Grant:          &access.Grant{},
					ApprovalMethod: &reviewed,
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
					Status:         access.APPROVED, // request should be approved
					RequestedBy:    "b",
					UpdatedAt:      clk.Now(),
					Grant:          &access.Grant{},
					ApprovalMethod: &reviewed,
				},
			},
		},
	}

	for _, tc := range testcases {

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			workflowMock := accessMocks.NewMockWorkflow(ctrl)
			if tc.wantErr == nil {
				workflowMock.EXPECT().Grant(gomock.Any(), gomock.Any(), gomock.Any()).Return(tc.withCreateGrantResponse.request.Grant, tc.withCreateGrantResponse.err).AnyTimes()
			}

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
				EventPutter: ep,
				Workflow:    workflowMock,
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
