package storage

import (
	"testing"

	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/ddb"
	"github.com/common-fate/ddb/ddbtest"
)

func TestGetRequestGroupWithTargetsForReviewer(t *testing.T) {
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

	ddbtest.PutFixtures(t, ts.db, []ddb.Keyer{&req, &group, &target, &req2, &group2, &target2})

	tc := []ddbtest.QueryTestCase{
		{
			Name: "ok",
			Query: &GetRequestGroupWithTargetsForReviewer{
				ReviewerID: "rev_abcd",
				RequestID:  req.ID,
				GroupID:    group.ID,
			},
			Want: &GetRequestGroupWithTargetsForReviewer{
				ReviewerID: "rev_abcd",
				RequestID:  req.ID,
				GroupID:    group.ID,
				Result: &access.GroupWithTargets{
					Group:   group,
					Targets: []access.GroupTarget{target},
				},
			},
		},
		{
			Name: "No access",
			Query: &GetRequestGroupWithTargetsForReviewer{
				ReviewerID: "rev_abcd",
				GroupID:    gid,
				RequestID:  req2.ID,
			},
			WantErr: ddb.ErrNoItems,
		},
	}

	ddbtest.RunQueryTests(t, ts.db, tc)
}
