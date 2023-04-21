package storage

import (
	"testing"

	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/ddb"
	"github.com/common-fate/ddb/ddbtest"
)

func TestGetRequestWithGroupsWithTargets(t *testing.T) {
	db := newTestingStorage(t)
	rid := "req_abcd"
	gid := "grp_abcd"
	tid := "gta_abcd"
	req := access.Request{ID: rid}
	group := access.Group{ID: gid, RequestID: rid}
	target := access.GroupTarget{ID: tid, GroupID: gid, RequestID: rid}

	ddbtest.PutFixtures(t, db, []ddb.Keyer{&req, &group, &target})

	tc := []ddbtest.QueryTestCase{
		{
			Name:  "ok",
			Query: &GetRequestWithGroupsWithTargets{ID: rid},
			Want: &GetRequestWithGroupsWithTargets{ID: rid, Result: &access.RequestWithGroupsWithTargets{
				Request: req,
				Groups: []access.GroupWithTargets{{
					Group:   group,
					Targets: []access.GroupTarget{target},
				}},
			}},
		},
		{
			Name:    "target group not found",
			Query:   &GetRequestWithGroupsWithTargets{ID: "not-found"},
			WantErr: ddb.ErrNoItems,
		},
	}

	ddbtest.RunQueryTests(t, db, tc)
}
