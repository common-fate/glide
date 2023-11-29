package slacknotifier

import (
	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExtractSlackMentions(t *testing.T) {
	type testcase struct {
		name     string
		rule     rule.AccessRule
		expected []slackMention
	}

	testCases := []testcase{
		{
			name:     "No matches",
			rule:     rule.TestAccessRule(rule.WithDescription("Here is my generic description.")),
			expected: make([]slackMention, 0),
		},
		{
			name:     "Matches only slack channel",
			rule:     rule.TestAccessRule(rule.WithDescription("Here is my generic description.\n This rule notifies @slack-C1234")),
			expected: []slackMention{{Channel: "C1234"}},
		},
		{
			name:     "Has slack channel and group",
			rule:     rule.TestAccessRule(rule.WithDescription("Here is my generic description.\n This rule notifies @slack-C1234 <@foo>")),
			expected: []slackMention{{Channel: "C1234", User: "@foo"}},
		},
		{
			name:     "Allows hyphens in channel name",
			rule:     rule.TestAccessRule(rule.WithDescription("Here is my generic description.\n This rule notifies @slack-C12-34 <@foo>")),
			expected: []slackMention{{Channel: "C12-34", User: "@foo"}},
		},
		{
			name: "Allows multiple channels and mentions.",
			rule: rule.TestAccessRule(rule.WithDescription("Here is my generic description.\n This rule notifies @slack-CABC <@foo>\n @slack-CXYZ <@bar>")),
			expected: []slackMention{
				{Channel: "CABC", User: "@foo"},
				{Channel: "CXYZ", User: "@bar"},
			},
		},
		{
			name: "Allows mix of channels and mentions.",
			rule: rule.TestAccessRule(rule.WithDescription("Here is my generic description.\n This rule notifies @slack-CABC @slack-CXYZ <@bar>")),
			expected: []slackMention{
				{Channel: "CABC"},
				{Channel: "CXYZ", User: "@bar"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.expected, extractSlackMentions(&tc.rule))
		})
	}
}
