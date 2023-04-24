package storage

import (
	"testing"

	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/ddb"
	"github.com/common-fate/ddb/ddbtest"
)

func TestListRequestWithGroupsWithTargetsForReviewer(t *testing.T) {
	ts := newTestingStorage(t)
	rid := "req_abcd"
	gid := "grp_abcd"
	tid := "gta_abcd"
	req := access.Request{ID: rid, GroupTargetCount: 1, RequestReviewers: []string{"rev_abcd"}}
	group := access.Group{ID: gid, RequestID: rid, RequestReviewers: []string{"rev_abcd"}}
	target := access.GroupTarget{ID: tid, GroupID: gid, RequestID: rid, RequestReviewers: []string{"rev_abcd"}}
	rid = "req_efgh"
	gid = "grp_efgh"
	tid = "gta_efgh"
	req2 := access.Request{ID: rid, GroupTargetCount: 1}
	group2 := access.Group{ID: gid, RequestID: rid}
	target2 := access.GroupTarget{ID: tid, GroupID: gid, RequestID: rid}
	// cleanup before the test
	err := ts.deleteAll()
	if err != nil {
		t.Fatal(err)
	}
	ddbtest.PutFixtures(t, ts.db, []ddb.Keyer{&req, &group, &target, &req2, &group2, &target2})

	tc := []ddbtest.QueryTestCase{
		{
			Name: "ok",
			Query: &ListRequestWithGroupsWithTargetsForReviewer{
				ReviewerID: "rev_abcd",
			},
			Want: &ListRequestWithGroupsWithTargetsForReviewer{
				ReviewerID: "rev_abcd",
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
			Query: &ListRequestWithGroupsWithTargetsForReviewer{
				ReviewerID: "rev_other",
			},
			Want: &ListRequestWithGroupsWithTargetsForReviewer{
				ReviewerID: "rev_other",
				Result:     []access.RequestWithGroupsWithTargets{},
			},
		},
	}

	ddbtest.RunQueryTests(t, ts.db, tc)
}
