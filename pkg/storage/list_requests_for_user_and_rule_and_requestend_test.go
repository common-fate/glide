package storage

import (
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/common-fate/ddb/ddbtest"
	"github.com/common-fate/granted-approvals/pkg/access"
	"github.com/common-fate/granted-approvals/pkg/types"
)

func TestListRequestsForUserAndRuleAndRequestend(t *testing.T) {
	s := newTestingStorage(t)

	user := types.NewUserID()
	rule := types.NewAccessRuleID()
	rule2 := types.NewAccessRuleID()
	clk := clock.NewMock()
	now := clk.Now()
	// set up fixture data for testing with.
	fixtureRequests := []access.Request{
		{
			ID:          types.NewRequestID(),
			RequestedBy: user,
			Rule:        rule,
			Status:      access.APPROVED,
			RequestedTiming: access.Timing{
				StartTime: &now,
				Duration:  time.Minute,
			},
		},
		{
			ID:          types.NewRequestID(),
			RequestedBy: user,
			Rule:        rule2,
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
			Query: &ListRequestsForUserAndRuleAndRequestend{
				UserID:               user,
				RuleID:               rule,
				RequestEndComparator: GreaterThanEqual,
				CompareTo:            now,
			},
			Want: &ListRequestsForUserAndRuleAndRequestend{
				UserID:               user,
				RuleID:               rule,
				RequestEndComparator: GreaterThanEqual,
				CompareTo:            now,
				Result:               fixtureRequests[:1],
			},
		},
		{
			Name: "ok2",
			Query: &ListRequestsForUserAndRuleAndRequestend{
				UserID:               user,
				RuleID:               rule,
				RequestEndComparator: LessThan,
				CompareTo:            now,
			},
			Want: &ListRequestsForUserAndRuleAndRequestend{
				UserID:               user,
				RuleID:               rule,
				RequestEndComparator: LessThan,
				CompareTo:            now,
				Result:               []access.Request{},
			},
		},
		{
			Name: "uther user",
			Query: &ListRequestsForUserAndRuleAndRequestend{
				UserID:               "other",
				RuleID:               rule,
				RequestEndComparator: GreaterThanEqual,
				CompareTo:            now,
			},
			Want: &ListRequestsForUserAndRuleAndRequestend{
				UserID:               "other",
				RuleID:               rule,
				RequestEndComparator: GreaterThanEqual,
				CompareTo:            now,
				Result:               []access.Request{},
			},
		},
	}

	ddbtest.RunQueryTests(t, s, tc)
}
