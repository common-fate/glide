package storage

import (
	"testing"

	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/ddb"
	"github.com/common-fate/ddb/ddbtest"
)

func TestGetRequestGroupWithTargets(t *testing.T) {
	ts := newTestingStorage(t)
	rid := "req_abcd"
	gid := "grp_abcd"
	tid := "gta_abcd"
	req := access.Request{ID: rid}
	group := access.Group{ID: gid, RequestID: rid}
	target := access.GroupTarget{ID: tid, GroupID: gid, RequestID: rid}

	ddbtest.PutFixtures(t, ts.db, []ddb.Keyer{&req, &group, &target})

	tc := []ddbtest.QueryTestCase{
		{
			Name:  "ok",
			Query: &GetRequestGroupWithTargets{RequestID: rid, GroupID: gid},
			Want: &GetRequestGroupWithTargets{RequestID: rid, GroupID: gid, Result: &access.GroupWithTargets{
				Group:   group,
				Targets: []access.GroupTarget{target},
			}},
		},
		{
			Name:    "target group not found",
			Query:   &GetRequestGroupWithTargets{RequestID: rid, GroupID: "not found"},
			WantErr: ddb.ErrNoItems,
		},
	}

	ddbtest.RunQueryTests(t, ts.db, tc)
}
