package storage

// import (
// 	"context"
// 	"testing"

// 	"github.com/common-fate/common-fate/pkg/targetgroup"
// 	"github.com/common-fate/ddb/ddbtest"
// 	"github.com/stretchr/testify/assert"
// )

// func TestGetHighestPriorityDeployment(t *testing.T) {

// 	type testcase struct {
// 		name                   string
// 		targetGroup            targetgroup.TargetGroup
// 		highPriorityDeployment targetgroup.Deployment
// 		lowPriorityDeployment  targetgroup.Deployment

// 		want    targetgroup.Deployment
// 		wantErr error
// 	}
// 	tg := targetgroup.TestTargetGroup()

// 	LowPriority := targetgroup.TestTargetGroupDeployment(func(d *targetgroup.Deployment) {
// 		d.TargetGroupAssignment = targetgroup.TargetGroupAssignment{TargetGroupID: tg.ID, Priority: 100, Valid: false}
// 		d.Healthy = true
// 	})
// 	HighPriority := targetgroup.TestTargetGroupDeployment(func(d *targetgroup.Deployment) {
// 		d.TargetGroupAssignment = targetgroup.TargetGroupAssignment{TargetGroupID: tg.ID, Priority: 999, Valid: true}
// 		d.Healthy = true
// 	})

// 	testcases := []testcase{
// 		{
// 			name:                   "ok",
// 			highPriorityDeployment: HighPriority,
// 			lowPriorityDeployment:  LowPriority,
// 			targetGroup:            tg,
// 			want:                   HighPriority,
// 		},
// 	}
// 	for _, tc := range testcases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			ctx := context.Background()
// 			db := newTestingStorage(t)

// 			ddbtest.PutFixtures(t, db, &tc.targetGroup)

// 			ddbtest.PutFixtures(t, db, &tc.highPriorityDeployment)
// 			ddbtest.PutFixtures(t, db, &tc.lowPriorityDeployment)

// 			q := GetTargetGroupDeploymentWithPriority{TargetGroupId: tc.targetGroup.ID}
// 			_, err := db.Query(ctx, &q)
// 			if err != nil && tc.wantErr == nil {
// 				t.Fatal(err)
// 			}
// 			got := q.Result

// 			if tc.wantErr != nil {
// 				assert.Equal(t, tc.wantErr, err)
// 			}

// 			assert.Contains(t, got.TargetGroupAssignment.Priority, tc.want.TargetGroupAssignment.Priority)

// 		})
// 	}

// }
