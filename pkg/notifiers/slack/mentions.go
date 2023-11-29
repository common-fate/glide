package slacknotifier

import (
	"github.com/common-fate/common-fate/pkg/rule"
	"regexp"
)

var slackRegex = regexp.MustCompile("@slack-(C[\\w-]+) ?(<(@[\\w]+)>)?")

// extractSlackMentions attempts to pull any slack tags from the rule description. Slack mentions are specified in the format
// of @slack-CHANNEL_ID @USER. The @USER tag is optional
func extractSlackMentions(rule *rule.AccessRule) []slackMention {
	matches := slackRegex.FindAllStringSubmatch(rule.Description, -1)

	mentions := make([]slackMention, 0)

	for _, groups := range matches {
		if len(groups) == 2 {
			mentions = append(mentions, slackMention{
				Channel: groups[1],
			})
		} else if len(groups) == 4 {
			mentions = append(mentions, slackMention{
				Channel: groups[1],
				User:    groups[3],
			})
		}
	}
	return mentions
}

type slackMention struct {
	Channel, User string
}
