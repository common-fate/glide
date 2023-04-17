package workflowsvc

import (
	"context"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/common-fate/common-fate/pkg/requests"
	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/common-fate/common-fate/pkg/service/workflowsvc/mocks"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb/ddbmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateGrant(t *testing.T) {
	type testcase struct {
		name                       string
		withCreateGrantResponseErr error
		subject                    string
		giveRequest                requests.AccessGroup
		wantErr                    error
		want                       []requests.Grantv2
	}
	clk := clock.NewMock()
	now := clk.Now()
	testcases := []testcase{
		{
			name: "ok",

			subject: "test@commonfate.io",
			giveRequest: requests.AccessGroup{
				AccessRule: rule.AccessRule{ID: "rule_a", Target: rule.Target{TargetGroupID: "test"}},
				ID:         "123",
				Request:    "abc",
				TimeConstraints: requests.Timing{
					Duration:  time.Hour,
					StartTime: &now,
				},
				With: []map[string]string{
					{
						"accountId":     "123",
						"permissionSet": "abc",
					},
					{
						"accountId":     "456",
						"permissionSet": "abc",
					},
				},
			},
			withCreateGrantResponseErr: nil,
			want: []requests.Grantv2{
				{
					ID:                 CreateGrantIdHash("test@commonfate.io", now, "test"),
					AccessGroup:        "123",
					Status:             types.GrantStatus(requests.PENDING),
					Start:              now,
					End:                now.Add(time.Hour),
					AccessInstructions: "",
					Subject:            "test@commonfate.io",
					CreatedAt:          now,
					UpdatedAt:          now,
					With: types.Grant_With{
						AdditionalProperties: map[string]string{
							"accountId":     "123",
							"permissionSet": "abc",
						},
					},
				},
				{
					ID:                 CreateGrantIdHash("test@commonfate.io", now, "test"),
					AccessGroup:        "123",
					Status:             types.GrantStatus(requests.PENDING),
					Start:              now,
					End:                now.Add(time.Hour),
					AccessInstructions: "",
					Subject:            "test@commonfate.io",
					CreatedAt:          now,
					UpdatedAt:          now,
					With: types.Grant_With{
						AdditionalProperties: map[string]string{
							"accountId":     "456",
							"permissionSet": "abc",
						},
					},
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			runtime := mocks.NewMockRuntime(ctrl)
			runtime.EXPECT().Grant(gomock.Any(), gomock.Any()).Return(tc.withCreateGrantResponseErr).AnyTimes()

			eventbus := mocks.NewMockEventPutter(ctrl)
			eventbus.EXPECT().Put(gomock.Any(), gomock.Any()).Return(tc.withCreateGrantResponseErr).AnyTimes()

			c := ddbmock.New(t)
			// c.MockQueryWithErr(&storage.GetUser{Result: tc.withUser}, tc.wantUserErr)

			s := Service{
				Runtime:  runtime,
				DB:       c,
				Clk:      clk,
				Eventbus: eventbus,
			}

			gotGrants, err := s.Grant(context.Background(), tc.giveRequest, tc.subject)
			assert.Equal(t, tc.wantErr, err)

			assert.Equal(t, tc.want, gotGrants)
		})
	}
}
