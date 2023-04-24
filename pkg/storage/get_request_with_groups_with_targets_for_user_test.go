package storage

import (
	"testing"

	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/ddb"
	"github.com/common-fate/ddb/ddbtest"
)

func TestGetRequestWithGroupsWithTargetsForUser(t *testing.T) {
	ts := newTestingStorage(t)
	rid := "req_abcd"
	gid := "grp_abcd"
	tid := "gta_abcd"
	req := access.Request{ID: rid, GroupTargetCount: 1, RequestedBy: "usr_abcd"}
	group := access.Group{ID: gid, RequestID: rid, RequestedBy: "usr_abcd"}
	target := access.GroupTarget{ID: tid, GroupID: gid, RequestID: rid, RequestedBy: "usr_abcd"}
	rid = "req_efgh"
	gid = "grp_efgh"
	tid = "gta_efgh"
	req2 := access.Request{ID: rid, GroupTargetCount: 1, RequestedBy: "usr_efgh"}
	group2 := access.Group{ID: gid, RequestID: rid, RequestedBy: "usr_efgh"}
	target2 := access.GroupTarget{ID: tid, GroupID: gid, RequestID: rid, RequestedBy: "usr_efgh"}

	ddbtest.PutFixtures(t, ts.db, []ddb.Keyer{&req, &group, &target, &req2, &group2, &target2})

	tc := []ddbtest.QueryTestCase{
		{
			Name: "ok",
			Query: &GetRequestWithGroupsWithTargetsForUser{
				UserID:    "usr_abcd",
				RequestID: req.ID,
			},
			Want: &GetRequestWithGroupsWithTargetsForUser{
				UserID:    "usr_abcd",
				RequestID: req.ID,
				Result: &access.RequestWithGroupsWithTargets{
					Request: req,
					Groups: []access.GroupWithTargets{{
						Group:   group,
						Targets: []access.GroupTarget{target},
					}},
				},
			},
		},
		{
			Name: "No access",
			Query: &GetRequestWithGroupsWithTargetsForUser{
				UserID:    "usr_abcd",
				RequestID: req2.ID,
			},
			WantErr: ddb.ErrNoItems,
		},
	}

	ddbtest.RunQueryTests(t, ts.db, tc)
}
