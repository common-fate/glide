package storage

import (
	"testing"

	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
	"github.com/common-fate/ddb/ddbtest"
)

func TestGetAccessRule(t *testing.T) {
	ts := newTestingStorage(t)

	rul := rule.TestAccessRule()
	ddbtest.PutFixtures(t, ts.db, &rul)

	tc := []ddbtest.QueryTestCase{
		{
			Name:  "ok",
			Query: &GetAccessRule{ID: rul.ID},
			Want:  &GetAccessRule{ID: rul.ID, Result: &rul},
		},
		{
			Name:    "rule not found",
			Query:   &GetAccessRule{ID: types.NewAccessRuleID()},
			WantErr: ddb.ErrNoItems,
		},
	}

	ddbtest.RunQueryTests(t, ts.db, tc)
}
