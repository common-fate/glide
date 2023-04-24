package storage

import (
	"testing"

	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
	"github.com/common-fate/ddb/ddbtest"
	"github.com/segmentio/ksuid"
)

func TestListUsers(t *testing.T) {
	ts := newTestingStorage(t)
	err := ts.deleteAll()
	if err != nil {
		t.Fatal(err)
	}
	active := identity.User{
		ID:     types.NewUserID(),
		Status: types.IdpStatusACTIVE,
	}

	archived := identity.User{
		ID:     ksuid.New().String(),
		Status: types.IdpStatusARCHIVED,
	}

	ddbtest.PutFixtures(t, ts.db, []ddb.Keyer{&active, &archived})
	tc := []ddbtest.QueryTestCase{
		{
			Name:  "ok",
			Query: &ListUsers{},
			Want:  &ListUsers{Result: []identity.User{active, archived}},
		},
	}

	ddbtest.RunQueryTests(t, ts.db, tc)
}
