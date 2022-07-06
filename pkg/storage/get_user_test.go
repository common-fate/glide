package storage

import (
	"testing"

	"github.com/common-fate/ddb"
	"github.com/common-fate/ddb/ddbtest"
	"github.com/common-fate/granted-approvals/pkg/identity"
	"github.com/common-fate/granted-approvals/pkg/types"
)

func TestGetUser(t *testing.T) {
	db := newTestingStorage(t)

	us := identity.User{
		ID: types.NewUserID(),
	}
	ddbtest.PutFixtures(t, db, &us)

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

	ddbtest.RunQueryTests(t, db, tc)
}
