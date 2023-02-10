package storage

import (
	"testing"

	"github.com/common-fate/common-fate/pkg/targetgroup"
	"github.com/common-fate/ddb"
	"github.com/common-fate/ddb/ddbtest"
)

func TestGetTargetGroupDeployment_BuildQuery(t *testing.T) {

	db := newTestingStorage(t)

	tg := targetgroup.TestTargetGroupDeployment("t1")
	ddbtest.PutFixtures(t, db, &tg)

	tc := []ddbtest.QueryTestCase{
		{
			Name:  "ok",
			Query: &GetTargetGroupDeployment{ID: tg.ID},
			Want:  &GetTargetGroupDeployment{ID: tg.ID, Result: tg},
		},
		{
			Name:    "target group deployment not found",
			Query:   &GetTargetGroupDeployment{ID: "not-found"},
			WantErr: ddb.ErrNoItems,
		},
	}

	ddbtest.RunQueryTests(t, db, tc)

}
