package accesssvc

import (
	"context"
	"testing"

	"github.com/benbjohnson/clock"
	"github.com/common-fate/ddb/ddbmock"
	"github.com/common-fate/granted-approvals/pkg/access"
	"github.com/common-fate/granted-approvals/pkg/storage"
	"github.com/stretchr/testify/assert"
)

func TestCancelRequest(t *testing.T) {
	type testcase struct {
		name               string
		givenCancelRequest CancelRequestOpts
		getRequestResponse *access.Request
		getRequestErr      error
		wantErr            error
	}

	clk := clock.NewMock()

	testcases := []testcase{
		{
			name: "ok",
			givenCancelRequest: CancelRequestOpts{
				CancellerID: "abcd",
				RequestID:   "req123",
			},
			getRequestResponse: &access.Request{
				RequestedBy: "abcd",
				Status:      access.PENDING,
			},
			wantErr: nil,
		},
		{
			name: "user not authorized",
			givenCancelRequest: CancelRequestOpts{
				CancellerID: "abcd",
				RequestID:   "req123",
			},
			getRequestResponse: &access.Request{
				RequestedBy: "other-user",
				Status:      access.PENDING,
			},
			wantErr: ErrUserNotAuthorized,
		},
		{
			name: "request not pending",
			givenCancelRequest: CancelRequestOpts{
				CancellerID: "abcd",
				RequestID:   "req123",
			},
			getRequestResponse: &access.Request{
				RequestedBy: "abcd",
				Status:      access.APPROVED,
			},
			wantErr: ErrRequestCannotBeCancelled,
		},
		{
			name: "unauthorised preceeds cannot be cancelled",
			givenCancelRequest: CancelRequestOpts{
				CancellerID: "abcd",
				RequestID:   "req123",
			},
			getRequestResponse: &access.Request{
				RequestedBy: "other-user",
				Status:      access.APPROVED,
			},
			wantErr: ErrUserNotAuthorized,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			db := ddbmock.New(t)
			db.MockQueryWithErr(&storage.GetRequest{Result: tc.getRequestResponse}, tc.getRequestErr)
			db.MockQuery(&storage.ListRequestReviewers{Result: []access.Reviewer{}})
			s := Service{
				Clock: clk,
				DB:    db,
			}
			err := s.CancelRequest(context.Background(), tc.givenCancelRequest)
			assert.Equal(t, tc.wantErr, err)

		})
	}

}
