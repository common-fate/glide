package storage

import (
	"testing"

	"github.com/common-fate/common-fate/pkg/handler"
	"github.com/common-fate/ddb"
	"github.com/common-fate/ddb/ddbtest"
)

func TestListHandlers(t *testing.T) {
	ts := newTestingStorage(t)
	err := ts.deleteAll()
	if err != nil {
		t.Fatal(err)
	}
	h := handler.TestHandler("test")

	ddbtest.PutFixtures(t, ts.db, []ddb.Keyer{&h})

	tc := []ddbtest.QueryTestCase{
		{
			Name:  "ok",
			Query: &ListHandlers{},
			Want:  &ListHandlers{Result: []handler.Handler{h}},
		},
	}

	ddbtest.RunQueryTests(t, ts.db, tc)
}
