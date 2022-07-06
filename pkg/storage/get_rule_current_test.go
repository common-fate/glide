package storage

import (
	"testing"

	"github.com/common-fate/ddb"
	"github.com/common-fate/ddb/ddbtest"
	"github.com/common-fate/granted-approvals/pkg/rule"
	"github.com/common-fate/granted-approvals/pkg/types"
)

func TestGetAccessRuleCurrent(t *testing.T) {
	db := newTestingStorage(t)

	rul := rule.TestAccessRule()
	ddbtest.PutFixtures(t, db, &rul)

	tc := []ddbtest.QueryTestCase{
		{
			Name:  "ok",
			Query: &GetAccessRuleCurrent{ID: rul.ID},
			Want:  &GetAccessRuleCurrent{ID: rul.ID, Result: &rul},
		},
		{
			Name:    "rule not found",
			Query:   &GetAccessRuleCurrent{ID: types.NewAccessRuleID()},
			WantErr: ddb.ErrNoItems,
		},
	}

	ddbtest.RunQueryTests(t, db, tc)
}
