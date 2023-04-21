package storage

import (
	"testing"

	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
	"github.com/common-fate/ddb/ddbtest"
	"github.com/segmentio/ksuid"
)

func TestGetGroup(t *testing.T) {
	db := newTestingStorage(t)

	g := identity.Group{
		ID:     ksuid.New().String(),
		Name:   "a",
		IdpID:  "a",
		Status: types.IdpStatusACTIVE,
	}
	ddbtest.PutFixtures(t, db, &g)

	tc := []ddbtest.QueryTestCase{
		{
			Name:  "ok",
			Query: &GetGroup{ID: g.ID},
			Want:  &GetGroup{ID: g.ID, Result: &g},
		},
		{
			Name:    "user not found",
			Query:   &GetGroup{ID: ksuid.New().String()},
			WantErr: ddb.ErrNoItems,
		},
	}

	ddbtest.RunQueryTests(t, db, tc)
}
