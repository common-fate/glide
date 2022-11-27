package accesssvc

import (
	"context"
	"testing"

	"github.com/benbjohnson/clock"
	"github.com/common-fate/common-fate/accesshandler/pkg/types"
	"github.com/common-fate/common-fate/pkg/access"
	accessMocks "github.com/common-fate/common-fate/pkg/service/accesssvc/mocks"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/ddb/ddbmock"
	"github.com/golang/mock/gomock"
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
			name: "active request not pending",
			givenCancelRequest: CancelRequestOpts{
				CancellerID: "abcd",
				RequestID:   "req123",
			},
			getRequestResponse: &access.Request{
				RequestedBy: "abcd",
				Status:      access.APPROVED,
				Grant:       &access.Grant{Status: types.GrantStatusACTIVE},
			},
			wantErr: ErrRequestCannotBeCancelled,
		},
		{
			name: "failed grant auto approved request can be cancelled",
			givenCancelRequest: CancelRequestOpts{
				CancellerID: "abcd",
				RequestID:   "req123",
			},
			getRequestResponse: &access.Request{
				RequestedBy: "abcd",
				Status:      access.APPROVED,
			},
			wantErr: nil,
		},
		{
			name: "cancelled request cannot be cancelled",
			givenCancelRequest: CancelRequestOpts{
				CancellerID: "abcd",
				RequestID:   "req123",
			},
			getRequestResponse: &access.Request{
				RequestedBy: "abcd",
				Status:      access.CANCELLED,
			},
			wantErr: ErrRequestCannotBeCancelled,
		},
		{
			name: "revoked grants cannot be cancelled",
			givenCancelRequest: CancelRequestOpts{
				CancellerID: "abcd",
				RequestID:   "req123",
			},
			getRequestResponse: &access.Request{
				RequestedBy: "abcd",
				Status:      access.APPROVED,
				Grant:       &access.Grant{Status: types.GrantStatusREVOKED},
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

			ctrl2 := gomock.NewController(t)
			ep := accessMocks.NewMockEventPutter(ctrl2)
			ep.EXPECT().Put(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

			s := Service{
				Clock:       clk,
				DB:          db,
				EventPutter: ep,
			}
			err := s.CancelRequest(context.Background(), tc.givenCancelRequest)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
