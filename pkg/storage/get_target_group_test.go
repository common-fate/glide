package storage

import (
	"testing"

	"github.com/common-fate/common-fate/pkg/target"
	"github.com/common-fate/ddb"
	"github.com/common-fate/ddb/ddbtest"
)

func TestGetTargetGroup(t *testing.T) {
	ts := newTestingStorage(t)

	tg := target.TestGroup()
	ddbtest.PutFixtures(t, ts.db, &tg)

	tc := []ddbtest.QueryTestCase{
		{
			Name:  "ok",
			Query: &GetTargetGroup{ID: tg.ID},
			Want:  &GetTargetGroup{ID: tg.ID, Result: &tg},
		},
		{
			Name:    "target group not found",
			Query:   &GetTargetGroup{ID: "not-found"},
			WantErr: ddb.ErrNoItems,
		},
	}

	ddbtest.RunQueryTests(t, ts.db, tc)
}
