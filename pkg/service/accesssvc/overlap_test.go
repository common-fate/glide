package accesssvc

import (
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/common-fate/granted-approvals/pkg/access"
	"github.com/common-fate/granted-approvals/pkg/rule"
	"github.com/common-fate/granted-approvals/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestOverlapsExistingGrant(t *testing.T) {
	type testcase struct {
		name               string
		accessRequest      access.Request
		upcomingRequests   []access.Request
		currentRequestRule rule.AccessRule
		allRules           []rule.AccessRule
		want               bool
		clock              clock.Clock
	}
	clk := clock.NewMock()
	inOneMinute := clk.Now().Add(time.Minute)

	now := clk.Now()

	//some requests premade with specific timings
	activeRequest := access.Request{Status: access.Status(types.RequestStatusAPPROVED), Grant: &access.Grant{Status: "ACTIVE"}, Rule: "rule_a", RequestedTiming: access.Timing{StartTime: &now, Duration: time.Minute * 2}}
	cancelledRequest := access.Request{Status: access.Status(types.RequestStatusCANCELLED), RequestedTiming: access.Timing{StartTime: &now, Duration: time.Minute * 2}, Rule: "rule_a"}
	revokedRequest := access.Request{Status: access.Status(types.RequestStatusCANCELLED), Grant: &access.Grant{Status: "REVOKED"}, RequestedTiming: access.Timing{StartTime: &now, Duration: time.Minute * 2}, Rule: "rule_a"}

	args1 := make(map[string]string)
	args1["1"] = "arg1"
	args1["2"] = "arg2"

	args2 := make(map[string]string)
	args2["a"] = "argA"
	args2["b"] = "argB"

	testcases := []testcase{
		{
			name:               "no existing grants",
			accessRequest:      access.Request{},
			upcomingRequests:   []access.Request{},
			currentRequestRule: rule.AccessRule{Target: rule.Target{}},
			allRules:           []rule.AccessRule{},
			clock:              clk,
			want:               false,
		},
		{
			name:               "different provider passes",
			accessRequest:      access.Request{ID: "abc", Rule: "rule_b"},
			upcomingRequests:   []access.Request{{ID: "def", Rule: "rule_a"}},
			currentRequestRule: rule.AccessRule{ID: "rule_a", Target: rule.Target{ProviderID: "prov_a"}},
			allRules:           []rule.AccessRule{{ID: "rule_a", Target: rule.Target{ProviderID: "prov_b"}}},
			clock:              clk,
			want:               false,
		},
		{
			name:               "request overlaps current active request fails",
			accessRequest:      access.Request{Rule: "rule_a", RequestedTiming: access.Timing{StartTime: &now, Duration: time.Minute * 5}},
			upcomingRequests:   []access.Request{activeRequest},
			currentRequestRule: rule.AccessRule{ID: "rule_a", Target: rule.Target{ProviderID: "prov_a"}},
			allRules:           []rule.AccessRule{{ID: "rule_a", Target: rule.Target{ProviderID: "prov_a"}}},
			clock:              clk,
			want:               true,
		},
		{
			name:               "scheduled request overlaps current active request fails",
			accessRequest:      access.Request{Rule: "rule_a", RequestedTiming: access.Timing{StartTime: &inOneMinute, Duration: time.Minute * 5}},
			upcomingRequests:   []access.Request{activeRequest},
			currentRequestRule: rule.AccessRule{ID: "rule_a", Target: rule.Target{ProviderID: "prov_a"}},
			allRules:           []rule.AccessRule{{ID: "rule_a", Target: rule.Target{ProviderID: "prov_a"}}},
			clock:              clk,
			want:               true,
		},

		{
			name:               "same rule different arguments should succeed",
			accessRequest:      access.Request{Rule: "rule_a", RequestedTiming: access.Timing{StartTime: &inOneMinute, Duration: time.Minute * 5}},
			upcomingRequests:   []access.Request{activeRequest},
			currentRequestRule: rule.AccessRule{ID: "rule_a", Target: rule.Target{ProviderID: "prov_a", With: args1}},
			allRules:           []rule.AccessRule{{ID: "rule_a", Target: rule.Target{ProviderID: "prov_a", With: args2}}},
			clock:              clk,
			want:               false,
		},
		{
			name:               "same rule same arguments should fail",
			accessRequest:      access.Request{Rule: "rule_a", RequestedTiming: access.Timing{StartTime: &inOneMinute, Duration: time.Minute * 5}},
			upcomingRequests:   []access.Request{activeRequest},
			currentRequestRule: rule.AccessRule{ID: "rule_a", Target: rule.Target{ProviderID: "prov_a", With: args1}},
			allRules:           []rule.AccessRule{{ID: "rule_a", Target: rule.Target{ProviderID: "prov_a", With: args1}}},
			clock:              clk,
			want:               true,
		},
		{
			name:               "same rule same arguments on expired request should pass",
			accessRequest:      access.Request{Rule: "rule_a", RequestedTiming: access.Timing{StartTime: &now, Duration: time.Minute * 5}},
			upcomingRequests:   []access.Request{cancelledRequest},
			currentRequestRule: rule.AccessRule{ID: "rule_a", Target: rule.Target{ProviderID: "prov_a", With: args1}},
			allRules:           []rule.AccessRule{{ID: "rule_a", Target: rule.Target{ProviderID: "prov_a", With: args1}}},
			clock:              clk,
			want:               false,
		},
		{
			name:               "same rule same arguments on revoked request should pass",
			accessRequest:      access.Request{Rule: "rule_a", RequestedTiming: access.Timing{StartTime: &now, Duration: time.Minute * 5}},
			upcomingRequests:   []access.Request{revokedRequest},
			currentRequestRule: rule.AccessRule{ID: "rule_a", Target: rule.Target{ProviderID: "prov_a", With: args1}},
			allRules:           []rule.AccessRule{{ID: "rule_a", Target: rule.Target{ProviderID: "prov_a", With: args1}}},
			clock:              clk,
			want:               false,
		},
		{
			name:          "overlap where one rule has selectable with, one rule has regular with",
			accessRequest: access.Request{Rule: "rule_a", RequestedTiming: access.Timing{StartTime: &now, Duration: time.Minute * 5}},
			upcomingRequests: []access.Request{{
				Rule: "rule_b",
				Grant: &access.Grant{
					Status: "ACTIVE",
				},
				RequestedTiming: access.Timing{StartTime: &now, Duration: time.Minute * 3},
				SelectedWith: map[string]access.Option{
					"a": {Value: "a"},
				},
			}},
			currentRequestRule: rule.AccessRule{
				ID: "rule_a",
				Target: rule.Target{
					ProviderID: "prov_a",
					With: map[string]string{
						"a": "a",
					},
				},
			},
			allRules: []rule.AccessRule{
				{
					ID: "rule_a",
					Target: rule.Target{
						ProviderID: "prov_a",
						With: map[string]string{
							"a": "a",
						},
					},
				},
				{
					ID: "rule_b",
					Target: rule.Target{
						ProviderID: "prov_a",
						WithSelectable: map[string][]string{
							"a": {"a"},
						},
					},
				},
			},
			clock: clk,
			want:  true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := overlapsExistingGrantCheck(tc.accessRequest, tc.upcomingRequests, tc.currentRequestRule, tc.allRules, tc.clock)
			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})

	}

}
