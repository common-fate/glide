package storage

import (
	"context"
	"testing"

	"github.com/common-fate/common-fate/pkg/target"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb/ddbtest"
	"github.com/stretchr/testify/assert"
)

// This test asserts that ordering is from highest to lowest and that only valid routes are returned
func TestListValidTargetRoutesForGroupByPriority(t *testing.T) {
	db := newTestingStorage(t)
	groupID := types.NewGroupID()
	r1 := target.Route{
		Group:    groupID,
		Handler:  types.NewGroupID(),
		Mode:     "Default",
		Priority: 1,
		Valid:    true,
	}
	r2 := target.Route{
		Group:    groupID,
		Handler:  types.NewGroupID(),
		Mode:     "Default",
		Priority: 100,
		Valid:    true,
	}
	r3 := target.Route{
		Group:    groupID,
		Handler:  types.NewGroupID(),
		Mode:     "Default",
		Priority: 200,
		Valid:    false,
	}
	ddbtest.PutFixtures(t, db, []target.Route{r1, r2, r3})
	q := &ListValidTargetRoutesForGroupByPriority{
		Group: groupID,
	}
	_, err := db.Query(context.TODO(), q)
	assert.NoError(t, err)
	assert.Equal(t, []target.Route{r2, r1}, q.Result)
}
