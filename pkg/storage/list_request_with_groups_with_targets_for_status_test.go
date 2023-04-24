package storage

import (
	"context"
	"testing"

	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
	"github.com/common-fate/ddb/ddbtest"
)

func TestListRequestWithGroupsWithTargetsForStatus(t *testing.T) {
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
	req2 := access.Request{ID: rid, GroupTargetCount: 1, RequestedBy: "usr_efgh", RequestStatus: types.ACTIVE}
	group2 := access.Group{ID: gid, RequestID: rid, RequestedBy: "usr_efgh", RequestStatus: types.ACTIVE}
	target2 := access.GroupTarget{ID: tid, GroupID: gid, RequestID: rid, RequestedBy: "usr_efgh", RequestStatus: types.ACTIVE}

	// cleanup before the test
	err := deleteAllRequests(context.Background(), db)
	if err != nil {
		t.Fatal(err)
	}
	ddbtest.PutFixtures(t, db, []ddb.Keyer{&req, &group, &target, &req2, &group2, &target2})

	// this test asserts that items are retrieved correctly and in the expected order, most recently created upcoming request to oldest created past request
	testcases := []ddbtest.QueryTestCase{
		{
			Name: "ok complete",
			Query: &ListRequestWithGroupsWithTargetsForStatus{
				Status: types.COMPLETE,
			},
			Want: &ListRequestWithGroupsWithTargetsForStatus{
				Status: types.COMPLETE,
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
			Query: &ListRequestWithGroupsWithTargetsForStatus{
				Status: types.CANCELLED,
			},
			Want: &ListRequestWithGroupsWithTargetsForStatus{
				Status: types.CANCELLED,
			},
		},
	}

	ddbtest.RunQueryTests(t, db, testcases)

}
