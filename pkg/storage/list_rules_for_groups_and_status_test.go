package storage

import (
	"testing"

	"github.com/common-fate/ddb"
	"github.com/common-fate/ddb/ddbtest"
	"github.com/common-fate/granted-approvals/pkg/rule"
	"github.com/segmentio/ksuid"
)

func TestListAccessRulesForGroupsAndStatus(t *testing.T) {
	s := newTestingStorage(t)

	// set up fixture data for testing with.

	// first fixture is a rule with a random group ID.
	group1 := ksuid.New().String()
	rule1 := rule.TestAccessRule(rule.WithName("rule1"), rule.WithGroups(group1))
	ddbtest.PutFixtures(t, s, &rule1)

	// rule2 has two groups associated with it
	group2a := ksuid.New().String()
	group2b := ksuid.New().String()
	rule2 := rule.TestAccessRule(rule.WithName("rule2"), rule.WithGroups(group2a, group2b))
	ddbtest.PutFixtures(t, s, &rule2)

	// rule3 has no groups associated with it.
	rule3 := rule.TestAccessRule(rule.WithName("rule3"), rule.WithGroups())
	ddbtest.PutFixtures(t, s, &rule3)

	//rule 4 is archived
	rule4 := rule.TestAccessRule(rule.WithName("rule4"), rule.WithStatus(rule.ARCHIVED), rule.WithGroups(group1))
	ddbtest.PutFixtures(t, s, &rule4)

	tc := []ddbtest.QueryTestCase{
		{
			Name:  "ok",
			Query: &ListAccessRulesForGroupsAndStatus{Status: rule.ACTIVE, Groups: []string{group1}},
			Want:  &ListAccessRulesForGroupsAndStatus{Status: rule.ACTIVE, Groups: []string{group1}, Result: []rule.AccessRule{rule1}},
		},
		{
			Name:  "multiple groups",
			Query: &ListAccessRulesForGroupsAndStatus{Status: rule.ACTIVE, Groups: []string{group1, group2a}},
			Want:  &ListAccessRulesForGroupsAndStatus{Status: rule.ACTIVE, Groups: []string{group1, group2a}, Result: []rule.AccessRule{rule1, rule2}},
		},
		{
			Name:    "no groups",
			Query:   &ListAccessRulesForGroupsAndStatus{Status: rule.ACTIVE, Groups: []string{}},
			WantErr: ddb.ErrNoItems,
		},
		{
			Name:  "archived",
			Query: &ListAccessRulesForGroupsAndStatus{Status: rule.ARCHIVED, Groups: []string{group1}},
			Want:  &ListAccessRulesForGroupsAndStatus{Status: rule.ARCHIVED, Groups: []string{group1}, Result: []rule.AccessRule{rule4}},
		},
	}

	ddbtest.RunQueryTests(t, s, tc)
}
