package storage

import (
	"context"
	"testing"

	"github.com/common-fate/common-fate/pkg/target"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb/ddbtest"
	"github.com/stretchr/testify/assert"
)

func TestListTargetRoutesForHandler(t *testing.T) {
	db := newTestingStorage(t)

	tr1 := target.Route{
		Group:   types.NewGroupID(),
		Handler: types.NewGroupID(),
	}
	ddbtest.PutFixtures(t, db, &tr1)
	q := &ListTargetRoutesForHandler{
		Handler: tr1.Handler,
	}
	_, err := db.Query(context.TODO(), q)
	assert.NoError(t, err)
	assert.Contains(t, q.Result, tr1)
}
