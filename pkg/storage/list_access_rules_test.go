package storage

import (
	"testing"

	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/common-fate/ddb"
	"github.com/common-fate/ddb/ddbtest"
)

func TestListAccessRules(t *testing.T) {
	ts := newTestingStorage(t)
	err := ts.deleteAll()
	if err != nil {
		t.Fatal(err)
	}

	a := rule.TestAccessRule()
	a.Priority = 100
	b := rule.TestAccessRule()
	b.Priority = 0

	ddbtest.PutFixtures(t, ts.db, []ddb.Keyer{&a, &b})

	tc := []ddbtest.QueryTestCase{
		{
			Name:  "ok",
			Query: &ListAccessRules{},
			// asserts the order is from highest to lowest priority
			Want: &ListAccessRules{Result: []rule.AccessRule{a, b}},
		},
	}

	ddbtest.RunQueryTests(t, ts.db, tc, ddbtest.WithAssertResultsOrder(true))

}
