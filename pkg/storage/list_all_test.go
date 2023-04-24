package storage

import (
	"context"
	"testing"

	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"

	"github.com/common-fate/ddb/ddbtest"
	"github.com/stretchr/testify/assert"
)

// This is a test that asserts that ddb behaves as expected, it should be shifted to the ddb package
func TestListAll(t *testing.T) {
	ts := newTestingStorage(t)
	err := ts.deleteAll()
	if err != nil {
		t.Fatal(err)
	}
	g1 := identity.Group{
		ID:     "id1",
		Name:   "id1",
		IdpID:  "id1",
		Users:  []string{"a"},
		Status: types.IdpStatusACTIVE,
	}
	g2 := identity.Group{
		ID:     "id2",
		Name:   "id2",
		IdpID:  "id2",
		Users:  []string{"a"},
		Status: types.IdpStatusACTIVE,
	}
	g3 := identity.Group{
		ID:     "id3",
		Name:   "id3",
		IdpID:  "id3",
		Users:  []string{"a"},
		Status: types.IdpStatusACTIVE,
	}
	ddbtest.PutFixtures(t, ts.db, []identity.Group{g1, g2, g3})
	q := &ListGroups{}

	// in this example, limit applies a limit per query, so all results will be fetched at a rate of 1 result at a time
	err = ts.db.All(context.Background(), q, ddb.Limit(2))
	assert.NoError(t, err)
	assert.Equal(t, []identity.Group{g1, g2, g3}, q.Result)
}
