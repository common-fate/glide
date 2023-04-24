package storage

import (
	"testing"

	"github.com/common-fate/common-fate/pkg/target"
	"github.com/common-fate/ddb"
	"github.com/common-fate/ddb/ddbtest"
)

func TestListTargetGroups(t *testing.T) {
	ts := newTestingStorage(t)
	err := ts.deleteAll()
	if err != nil {
		t.Fatal(err)
	}

	tg1 := target.TestGroup()
	tg2 := target.TestGroup()
	ddbtest.PutFixtures(t, ts.db, []ddb.Keyer{&tg1, &tg2})
	tc := []ddbtest.QueryTestCase{
		{
			Name:  "ok",
			Query: &ListTargetGroups{},
			Want:  &ListTargetGroups{Result: []target.Group{tg1, tg2}},
		},
	}

	ddbtest.RunQueryTests(t, ts.db, tc)
}
