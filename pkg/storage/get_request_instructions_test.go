package storage

import (
	"testing"

	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/ddb"
	"github.com/common-fate/ddb/ddbtest"
)

func TestGetRequestInstructions(t *testing.T) {
	db := newTestingStorage(t)

	i := access.Instructions{
		ID:           "test",
		Instructions: "example",
	}
	ddbtest.PutFixtures(t, db, &i)

	tc := []ddbtest.QueryTestCase{
		{
			Name:  "ok",
			Query: &GetRequestInstructions{ID: i.ID},
			Want:  &GetRequestInstructions{ID: i.ID, Result: &i},
		},
		{
			Name:    "instructions not found",
			Query:   &GetTargetGroup{ID: "not-found"},
			WantErr: ddb.ErrNoItems,
		},
	}

	ddbtest.RunQueryTests(t, db, tc)
}
