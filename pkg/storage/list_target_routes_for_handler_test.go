package storage

import (
	"testing"

	"github.com/common-fate/common-fate/pkg/target"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb/ddbtest"
)

func TestListTargetRoutesForHandler(t *testing.T) {
	ts := newTestingStorage(t)
	err := ts.deleteAll()
	if err != nil {
		t.Fatal(err)
	}
	tr1 := target.Route{
		Group:   types.NewGroupID(),
		Handler: types.NewGroupID(),
	}
	ddbtest.PutFixtures(t, ts.db, &tr1)

	tc := []ddbtest.QueryTestCase{
		{
			Name: "ok",
			Query: &ListTargetRoutesForHandler{
				Handler: tr1.Handler,
			},
			Want: &ListTargetRoutesForHandler{Handler: tr1.Handler, Result: []target.Route{tr1}},
		},
	}

	ddbtest.RunQueryTests(t, ts.db, tc)

}
