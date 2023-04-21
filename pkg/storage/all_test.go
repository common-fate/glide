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

func TestAll(t *testing.T) {
	db := newTestingStorage(t)

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
	ddbtest.PutFixtures(t, db, []identity.Group{g1, g2, g3})
	q := &ListGroups{}

	// in this example, limit applies a limit per query, so all results will be fetched at a rate of 1 result at a time
	err := db.All(context.Background(), q, ddb.Limit(1))

	assert.NoError(t, err)

	assert.Contains(t, q.Result, g1)
	assert.Contains(t, q.Result, g2)
	assert.Contains(t, q.Result, g3)
}
