package storage

import (
	"testing"

	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
	"github.com/common-fate/ddb/ddbtest"
)

func TestGetAccessRuleVersion(t *testing.T) {
	rul := rule.TestAccessRule()
	db := newTestingStorage(t)
	ddbtest.PutFixtures(t, db, &rul)

	tc := []ddbtest.QueryTestCase{
		{
			Name:  "ok",
			Query: &GetAccessRuleVersion{ID: rul.ID, VersionID: rul.Version},
			Want:  &GetAccessRuleVersion{ID: rul.ID, VersionID: rul.Version, Result: &rul},
		},
		{
			Name:    "rule not found",
			Query:   &GetAccessRuleVersion{ID: rul.ID, VersionID: types.NewVersionID()},
			WantErr: ddb.ErrNoItems,
		},
	}

	ddbtest.RunQueryTests(t, db, tc)

}
