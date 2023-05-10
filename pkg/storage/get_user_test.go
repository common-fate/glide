package storage

import (
	"testing"

	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
	"github.com/common-fate/ddb/ddbtest"
)

func TestGetUser(t *testing.T) {
	ts := newTestingStorage(t)

	us := identity.User{
		ID: types.NewUserID(),
	}
	ddbtest.PutFixtures(t, ts.db, &us)

	tc := []ddbtest.QueryTestCase{
		{
			Name:  "ok",
			Query: &GetUser{ID: us.ID},
			Want:  &GetUser{ID: us.ID, Result: &us},
		},
		{
			Name:    "user not found",
			Query:   &GetUser{ID: types.NewUserID()},
			WantErr: ddb.ErrNoItems,
		},
	}

	ddbtest.RunQueryTests(t, ts.db, tc)
}
