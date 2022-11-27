package storage

import (
	"testing"

	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb/ddbtest"
)

func TestListRequestEvents(t *testing.T) {
	s := newTestingStorage(t)

	reqID := types.NewRequestID()
	re1 := access.RequestEvent{ID: types.NewHistoryID(), RequestID: reqID}
	re2 := access.RequestEvent{ID: types.NewHistoryID(), RequestID: reqID}
	ddbtest.PutFixtures(t, s, []*access.RequestEvent{&re1, &re2})

	tc := []ddbtest.QueryTestCase{
		{
			Name:  "ok",
			Query: &ListRequestEvents{RequestID: reqID},
			Want:  &ListRequestEvents{RequestID: reqID, Result: []access.RequestEvent{re1, re2}},
		},
	}

	ddbtest.RunQueryTests(t, s, tc)
}
