package storage

import (
	"testing"

	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/ddb"
	"github.com/common-fate/ddb/ddbtest"
)

func TestListRequestWithGroupsWithTargetsForUser(t *testing.T) {
	db := newTestingStorage(t)
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

	ddbtest.PutFixtures(t, db, []ddb.Keyer{&req, &group, &target, &req2, &group2, &target2})

	tc := []ddbtest.QueryTestCase{
		{
			Name: "ok",
			Query: &ListRequestWithGroupsWithTargetsForUser{
				UserID: "usr_abcd",
			},
			Want: &ListRequestWithGroupsWithTargetsForUser{
				UserID: "usr_abcd",
				Result: []access.RequestWithGroupsWithTargets{
					{
						Request: req,
						Groups: []access.GroupWithTargets{{
							Group:   group,
							Targets: []access.GroupTarget{target},
						}},
					},
				},
			},
		},
		{
			Name: "No matches",
			Query: &ListRequestWithGroupsWithTargetsForUser{
				UserID: "usr_other",
			},
			Want: &ListRequestWithGroupsWithTargetsForUser{
				UserID: "usr_other",
				Result: []access.RequestWithGroupsWithTargets{},
			},
		},
	}

	ddbtest.RunQueryTests(t, db, tc)
}
