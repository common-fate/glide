package storage

import (
	"testing"

	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
	"github.com/common-fate/ddb/ddbtest"
)

func TestListGroups(t *testing.T) {

	ts := newTestingStorage(t)
	err := ts.deleteAll()
	if err != nil {
		t.Fatal(err)
	}

	active := identity.Group{
		ID:     types.NewGroupID(),
		Name:   "a",
		IdpID:  "a",
		Users:  []string{"a"},
		Status: types.IdpStatusACTIVE,
	}

	archived := identity.Group{
		ID:     types.NewGroupID(),
		Name:   "a",
		IdpID:  "a",
		Users:  []string{"a"},
		Status: types.IdpStatusARCHIVED,
	}

	ddbtest.PutFixtures(t, ts.db, []ddb.Keyer{&active, &archived})

	tc := []ddbtest.QueryTestCase{
		{
			Name:  "ok",
			Query: &ListGroups{},
			Want:  &ListGroups{Result: []identity.Group{active, archived}},
		},
	}

	ddbtest.RunQueryTests(t, ts.db, tc)
}
