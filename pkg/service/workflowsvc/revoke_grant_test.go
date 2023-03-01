package workflowsvc

import (
	"context"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/common-fate/common-fate/pkg/service/workflowsvc/mocks"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/ddb/ddbmock"
	"github.com/common-fate/iso8601"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type opts struct {
	Request   access.Request
	RevokerID string
}

func TestAccessRevoke(t *testing.T) {
	type testcase struct {
		name                        string
		give                        opts
		mockGetAccessRuleVersion    *rule.AccessRule
		mockGetAccessRuleVersionErr error
		wantErr                     error
		wantEventPutterErr          error
	}

	testStartTime := iso8601.Now().Add(time.Hour)
	testEndTime := iso8601.Now().Add(time.Hour * 2)

	testcases := []testcase{
		{
			name:    "Trying to revoke inactive grant",
			wantErr: ErrGrantInactive,
			give: opts{Request: access.Request{
				ID: "123",
				Grant: &access.Grant{
					Start:    testStartTime,
					End:      testEndTime,
					Subject:  "test@test.com",
					Status:   "PENDING",
					Provider: "okta",
				}}, RevokerID: "1234"}},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			runtime := mocks.NewMockRuntime(ctrl)
			runtime.EXPECT().Revoke(gomock.Any(), tc.give.Request.ID, gomock.Any()).Return(tc.wantErr).AnyTimes()

			eventPutter := mocks.NewMockEventPutter(ctrl)
			eventPutter.EXPECT().Put(gomock.Any(), gomock.Any()).Return(tc.wantEventPutterErr).AnyTimes()

			db := ddbmock.New(t)
			db.MockQueryWithErr(&storage.GetAccessRuleVersion{
				Result: tc.mockGetAccessRuleVersion,
			}, tc.mockGetAccessRuleVersionErr)
			s := Service{Runtime: runtime, Clk: clock.New(), DB: db, Eventbus: eventPutter}
			_, err := s.Revoke(context.Background(), tc.give.Request, tc.give.RevokerID)

			assert.Equal(t, tc.wantErr, err)
		})
	}

}
