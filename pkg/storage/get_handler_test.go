package storage

import (
	"testing"

	"github.com/common-fate/common-fate/pkg/handler"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
	"github.com/common-fate/ddb/ddbtest"
)

func TestGetHandler(t *testing.T) {
	db := newTestingStorage(t)

	hand := handler.TestHandler("test")
	ddbtest.PutFixtures(t, db, &hand)

	tc := []ddbtest.QueryTestCase{
		{
			Name:  "ok",
			Query: &GetHandler{ID: hand.ID},
			Want:  &GetHandler{ID: hand.ID, Result: &hand},
		},
		{
			Name:    "handler not found",
			Query:   &GetHandler{ID: types.NewUserID()},
			WantErr: ddb.ErrNoItems,
		},
	}

	ddbtest.RunQueryTests(t, db, tc)
}
