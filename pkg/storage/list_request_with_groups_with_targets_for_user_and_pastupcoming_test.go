package storage

import (
	"testing"

	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
	"github.com/common-fate/ddb/ddbtest"
)

func TestListRequestWithGroupsWithTargetsForUserAndPastUpcoming(t *testing.T) {
	db := newTestingStorage(t)
	rid := "req_abcd"
	gid := "grp_abcd"
	tid := "gta_abcd"
	req := access.Request{ID: rid, GroupTargetCount: 1, RequestedBy: "usr_abcd", RequestStatus: types.COMPLETE}
	group := access.Group{ID: gid, RequestID: rid, RequestedBy: "usr_abcd", RequestStatus: types.COMPLETE}
	target := access.GroupTarget{ID: tid, GroupID: gid, RequestID: rid, RequestedBy: "usr_abcd", RequestStatus: types.COMPLETE}
	rid = "req_efgh"
	gid = "grp_efgh"
	tid = "gta_efgh"
	req2 := access.Request{ID: rid, GroupTargetCount: 1, RequestedBy: "usr_abcd", RequestStatus: types.PENDING}
	group2 := access.Group{ID: gid, RequestID: rid, RequestedBy: "usr_abcd", RequestStatus: types.PENDING}
	target2 := access.GroupTarget{ID: tid, GroupID: gid, RequestID: rid, RequestedBy: "usr_abcd", RequestStatus: types.PENDING}

	ddbtest.PutFixtures(t, db, []ddb.Keyer{&req, &group, &target, &req2, &group2, &target2})

	tc := []ddbtest.QueryTestCase{
		{
			Name: "past",
			Query: &ListRequestWithGroupsWithTargetsForUserAndPastUpcoming{
				UserID:       "usr_abcd",
				PastUpcoming: keys.AccessRequestPastUpcomingPAST,
			},
			Want: &ListRequestWithGroupsWithTargetsForUserAndPastUpcoming{
				UserID:       "usr_abcd",
				PastUpcoming: keys.AccessRequestPastUpcomingPAST,
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
			Name: "upcoming",
			Query: &ListRequestWithGroupsWithTargetsForUserAndPastUpcoming{
				UserID:       "usr_abcd",
				PastUpcoming: keys.AccessRequestPastUpcomingUPCOMING,
			},
			Want: &ListRequestWithGroupsWithTargetsForUserAndPastUpcoming{
				UserID:       "usr_abcd",
				PastUpcoming: keys.AccessRequestPastUpcomingUPCOMING,
				Result: []access.RequestWithGroupsWithTargets{
					{
						Request: req2,
						Groups: []access.GroupWithTargets{{
							Group:   group2,
							Targets: []access.GroupTarget{target2},
						}},
					},
				},
			},
		},
		{
			Name: "No matches",
			Query: &ListRequestWithGroupsWithTargetsForUserAndPastUpcoming{
				UserID: "usr_other",
			},
			Want: &ListRequestWithGroupsWithTargetsForUserAndPastUpcoming{
				UserID: "usr_other",
				Result: []access.RequestWithGroupsWithTargets{},
			},
		},
	}

	ddbtest.RunQueryTests(t, db, tc)
}
