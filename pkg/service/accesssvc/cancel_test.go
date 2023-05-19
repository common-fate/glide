package accesssvc

import (
	"context"
	"testing"

	"github.com/benbjohnson/clock"
	"github.com/common-fate/common-fate/pkg/access"
	eventmock "github.com/common-fate/common-fate/pkg/eventhandler/mocks"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb/ddbmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCancelRequest(t *testing.T) {
	type testcase struct {
		name               string
		givenCancelRequest CancelRequestOpts
		getRequestResponse *access.RequestWithGroupsWithTargets
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
			getRequestResponse: &access.RequestWithGroupsWithTargets{
				Groups: []access.GroupWithTargets{
					{
						Group: access.Group{
							RequestID: "req123",
							Status:    types.RequestAccessGroupStatusPENDINGAPPROVAL,
						},
						Targets: []access.GroupTarget{},
					},
				},
				Request: access.Request{
					ID:            "req123",
					RequestStatus: types.RequestStatus(types.RequestAccessGroupStatusPENDINGAPPROVAL),
					RequestedBy: access.RequestedBy{
						ID: "abcd",
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "user not authorized",
			givenCancelRequest: CancelRequestOpts{
				CancellerID: "abcd",
				RequestID:   "req123",
			},
			getRequestResponse: &access.RequestWithGroupsWithTargets{
				Groups: []access.GroupWithTargets{
					{
						Group: access.Group{
							RequestID: "req123",
							Status:    types.RequestAccessGroupStatusPENDINGAPPROVAL,
						},
						Targets: []access.GroupTarget{},
					},
				},
				Request: access.Request{
					ID:            "req123",
					RequestStatus: types.RequestStatus(types.RequestAccessGroupStatusPENDINGAPPROVAL),
					RequestedBy: access.RequestedBy{
						ID: "different",
					},
				},
			},
			wantErr: ErrUserNotAuthorized,
		},
		{
			name: "active request not pending",
			givenCancelRequest: CancelRequestOpts{
				CancellerID: "abcd",
				RequestID:   "req123",
			},
			getRequestResponse: &access.RequestWithGroupsWithTargets{
				Groups: []access.GroupWithTargets{
					{
						Group: access.Group{
							RequestID: "req123",
							Status:    types.RequestAccessGroupStatusAPPROVED,
						},
						Targets: []access.GroupTarget{},
					},
				},
				Request: access.Request{
					ID:            "req123",
					RequestStatus: types.RequestStatus(types.RequestAccessGroupStatusPENDINGAPPROVAL),
					RequestedBy: access.RequestedBy{
						ID: "abcd",
					},
				},
			},
			wantErr: ErrRequestCannotBeCancelled,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			db := ddbmock.New(t)
			db.MockQueryWithErr(&storage.GetRequestWithGroupsWithTargets{Result: tc.getRequestResponse}, tc.getRequestErr)

			ctrl2 := gomock.NewController(t)
			ep := eventmock.NewMockEventPutter(ctrl2)
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
