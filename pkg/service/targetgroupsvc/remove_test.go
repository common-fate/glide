package targetgroupsvc

import (
	"context"
	"errors"
	"testing"

	"github.com/benbjohnson/clock"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/targetgroup"
	"github.com/common-fate/ddb/ddbmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRemoveTargetGroupLink(t *testing.T) {
	type testcase struct {
		name                                 string                  // ok
		giveDeploymentID                     string                  // ok
		giveGroupID                          string                  // ok
		mockGetTargetGroupResponse           targetgroup.TargetGroup // ok
		mockGetTargetGroupErr                error                   // ok
		mockGetTargetGroupDeploymentResponse targetgroup.Deployment  // ok
		mockGetTargetGroupDeploymentErr      error                   // ok
		// mockPutTargetGroupDeployment         *targetgroup.Deployment
		mockPutTargetGroupDeploymentErr error // ok
		want                            error // ok
	}

	testcases := []testcase{
		{
			name:                                 "ok",
			giveDeploymentID:                     "dep",
			giveGroupID:                          "grp",
			mockGetTargetGroupResponse:           targetgroup.TargetGroup{ID: "grp"},
			mockGetTargetGroupDeploymentResponse: targetgroup.Deployment{ID: "dep"},
			want:                                 nil,
		},
		{
			name:                  "target group err, error case",
			mockGetTargetGroupErr: errors.New("target group err"),
			want:                  errors.New("target group err"),
		},
		{
			name:                                 "target deployment error case",
			mockGetTargetGroupDeploymentResponse: targetgroup.Deployment{ID: "abc"},
			mockPutTargetGroupDeploymentErr:      errors.New("target deployment error"),
			want:                                 errors.New("target deployment error"),
		},
		{
			name:                                 "deployment update error case",
			mockGetTargetGroupDeploymentResponse: targetgroup.Deployment{ID: "abc"},
			mockGetTargetGroupErr:                errors.New("deployment update error"),
			want:                                 errors.New("deployment update error"),
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dbMock := ddbmock.New(t)

			// we need test cases 3 db queries
			// storage.GetTargetGroup
			// storage.GetTargetGroupDeployment
			// s.DB.Put
			// each following test case should have a mock for each of these queries

			if tc.mockGetTargetGroupErr != nil {
				dbMock.MockQueryWithErr(&storage.GetTargetGroup{ID: tc.giveGroupID}, tc.mockGetTargetGroupErr)
			} else {
				dbMock.MockQuery(&storage.GetTargetGroup{ID: tc.giveGroupID})
			}

			if tc.mockGetTargetGroupDeploymentErr != nil {
				dbMock.MockQueryWithErr(&storage.GetTargetGroupDeployment{ID: tc.giveDeploymentID}, tc.mockGetTargetGroupDeploymentErr)
			} else {
				dbMock.MockQuery(&storage.GetTargetGroupDeployment{ID: tc.giveDeploymentID})
			}

			if tc.mockPutTargetGroupDeploymentErr != nil {
				dbMock.PutErr = tc.mockPutTargetGroupDeploymentErr
			}

			clk := clock.NewMock()
			ctrl := gomock.NewController(t)

			defer ctrl.Finish()

			// new service
			s := Service{
				Clock: clk,
				DB:    dbMock,
			}

			err := s.RemoveTargetGroupLink(context.Background(), tc.giveDeploymentID, tc.giveGroupID)

			assert.Equal(t, tc.want, err)
		})
	}
}
