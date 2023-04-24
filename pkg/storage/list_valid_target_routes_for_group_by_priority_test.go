package storage

import (
	"testing"

	"github.com/common-fate/common-fate/pkg/target"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb/ddbtest"
)

// This test asserts that ordering is from highest to lowest and that only valid routes are returned
func TestListValidTargetRoutesForGroupByPriority(t *testing.T) {
	ts := newTestingStorage(t)
	groupID := types.NewGroupID()
	r1 := target.Route{
		Group:    groupID,
		Handler:  types.NewGroupID(),
		Kind:     "Default",
		Priority: 1,
		Valid:    true,
	}
	r2 := target.Route{
		Group:    groupID,
		Handler:  types.NewGroupID(),
		Kind:     "Default",
		Priority: 100,
		Valid:    true,
	}
	r3 := target.Route{
		Group:    groupID,
		Handler:  types.NewGroupID(),
		Kind:     "Default",
		Priority: 200,
		Valid:    false,
	}
	ddbtest.PutFixtures(t, ts.db, []target.Route{r1, r2, r3})

	tc := []ddbtest.QueryTestCase{
		{
			Name: "ok, invalid route is excluded",
			Query: &ListValidTargetRoutesForGroupByPriority{
				Group: groupID,
			},
			Want: &ListValidTargetRoutesForGroupByPriority{Group: groupID, Result: []target.Route{r2, r1}},
		},
	}

	ddbtest.RunQueryTests(t, ts.db, tc, ddbtest.WithAssertResultsOrder(true))
}
