package storage

import (
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb/ddbtest"
)

func TestListRequestsForUserAndRequestend(t *testing.T) {
	s := newTestingStorage(t)

	user := types.NewUserID()
	clk := clock.NewMock()
	now := clk.Now()
	// set up fixture data for testing with.
	fixtureRequests := []access.Request{
		{
			ID:          types.NewRequestID(),
			RequestedBy: user,
			Rule:        "randomRule",
			Status:      access.APPROVED,
			RequestedTiming: access.Timing{
				StartTime: &now,
				Duration:  time.Minute,
			},
		},
	}
	ddbtest.PutFixtures(t, s, fixtureRequests)

	tc := []ddbtest.QueryTestCase{
		{
			Name: "ok",
			Query: &ListRequestsForUserAndRequestend{
				UserID:               user,
				RequestEndComparator: GreaterThanEqual,
				CompareTo:            now,
			},
			Want: &ListRequestsForUserAndRequestend{
				UserID:               user,
				RequestEndComparator: GreaterThanEqual,
				CompareTo:            now,
				Result:               fixtureRequests,
			},
		},
		{
			Name: "ok",
			Query: &ListRequestsForUserAndRequestend{
				UserID:               user,
				RequestEndComparator: LessThan,
				CompareTo:            now,
			},
			Want: &ListRequestsForUserAndRequestend{
				UserID:               user,
				RequestEndComparator: LessThan,
				CompareTo:            now,
				Result:               []access.Request{},
			},
		},
		{
			Name: "uther user",
			Query: &ListRequestsForUserAndRequestend{
				UserID:               "other",
				RequestEndComparator: GreaterThanEqual,
				CompareTo:            now,
			},
			Want: &ListRequestsForUserAndRequestend{
				UserID:               "other",
				RequestEndComparator: GreaterThanEqual,
				CompareTo:            now,
				Result:               []access.Request{},
			},
		},
	}

	ddbtest.RunQueryTests(t, s, tc)
}
