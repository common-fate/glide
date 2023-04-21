package storage

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/ddb"
	"github.com/common-fate/ddb/ddbtest"
)

func TestListRequestWithGroupsWithTargets(t *testing.T) {
	db := newTestingStorage(t)
	rid := "req_abcd"
	gid := "grp_abcd"
	tid := "gta_abcd"
	req := access.Request{ID: rid, GroupTargetCount: 1}
	group := access.Group{ID: gid, RequestID: rid}
	target := access.GroupTarget{ID: tid, GroupID: gid, RequestID: rid}
	rid = "req_efgh"
	gid = "grp_efgh"
	tid = "gta_efgh"
	req2 := access.Request{ID: rid, GroupTargetCount: 1}
	group2 := access.Group{ID: gid, RequestID: rid}
	target2 := access.GroupTarget{ID: tid, GroupID: gid, RequestID: rid}

	// cleanup before the test
	err := deleteAllRequests(context.Background(), db)
	if err != nil {
		t.Fatal(err)
	}

	ddbtest.PutFixtures(t, db, []ddb.Keyer{&req, &group, &target, &req2, &group2, &target2})

	tc := []ddbtest.QueryTestCase{
		{
			Name:  "ok",
			Query: &ListRequestWithGroupsWithTargets{},
			Want: &ListRequestWithGroupsWithTargets{
				Result: []access.RequestWithGroupsWithTargets{
					{
						Request: req,
						Groups: []access.GroupWithTargets{{
							Group:   group,
							Targets: []access.GroupTarget{target},
						}},
					},
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
			QueryOpts: []func(*ddb.QueryOpts){ddb.Limit(4)},
			Name:      "paginated, returns only one complete request",
			Query:     &ListRequestWithGroupsWithTargets{},
			Want: &ListRequestWithGroupsWithTargets{
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
			// This asserts that the unmarshalling correctly set the pagination token to the key of the last target of teh last complete request which was unmarshalled
			WantNextPage: aws.String(`{"PK":{"Value":"ACCESS_REQUESTV2#"},"SK":{"Value":"ACCESS_REQUESTV2#req_abcd#ACCESS_REQUESTV2_GROUP#grp_abcd#ACCESS_REQUESTV2_GROUP_TARGET#gta_abcd#"}}`),
		},
		{
			QueryOpts: []func(*ddb.QueryOpts){ddb.Limit(2)},
			Name:      "paginated, returns only one complete request",
			Query:     &ListRequestWithGroupsWithTargets{},
			WantErr:   errors.New("failed to unmarshal requests, this could happen if the data for the request exceeds the 1mb limit for a ddb query"),
		},
	}

	ddbtest.RunQueryTests(t, db, tc)
}
