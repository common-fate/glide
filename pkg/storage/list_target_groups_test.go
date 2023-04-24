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

	tg1 := target.Group{
		ID:   "test-target-group1",
		Icon: "aws-sso",
		From: target.From{
			Publisher: "test",
			Name:      "test",
			Version:   "v1.1.1",
		},
	}
	tg2 := target.Group{
		ID:   "test-target-group2",
		Icon: "aws-sso",
		From: target.From{
			Publisher: "test",
			Name:      "test",
			Version:   "v1.1.1",
		},
	}
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
