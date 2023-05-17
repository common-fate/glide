package slacknotifier

import (
	"fmt"
	"strings"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/common-fate/common-fate/pkg/notifiers"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/slack-go/slack"
)

// titleCase turns STRING into String
func titleCase(s string) string {
	if s == "" {
		return ""
	}

	lower := strings.ToLower(s)
	cap := strings.ToUpper(string(lower[0])) + lower[1:]
	return cap
}

type RequestMessageOpts struct {
	Group access.Group
	// AccessGroups []access.Group
	// RequestArguments []types.With

	ReviewURLs notifiers.ReviewURLs
	// optional field that will replace the default requestor email with a slack @mention
	RequestorSlackID string
	RequestorEmail   string
	RequestReviewer  *identity.User
	WasReviewed      bool
	IsWebhook        bool
}

/*
*
BuildRequestReviewMessage builds a slack message for a request review based on the contextual RequestMessageOpts

Needs to handle:
Approved | Declined | Cancelled | Revoked

Goals:
- by using a fields for ReviewUrls etc. it allow us to easily mock the fn
*/
func BuildRequestReviewMessage(o RequestMessageOpts) (summary string, msg slack.Message) {

	group := o.Group

	requestor := group.RequestedBy.Email
	if o.RequestorSlackID != "" {
		requestor = fmt.Sprintf("<@%s>", o.RequestorSlackID)
	}

	status := string(group.Status)
	isCancelledOrRevoked := group.RequestStatus == types.CANCELLED || group.RequestStatus == types.REVOKED

	// In this scenario we want to override the request review status when,
	// The parent request has been: cancelled or revoked
	if isCancelledOrRevoked {
		status = string(o.Group.RequestStatus)
	}

	status = titleCase(string(status))
	statusLower := strings.ToLower(status)
	statusNoUnder := strings.Replace(statusLower, "_", " ", -1)

	if o.IsWebhook && o.WasReviewed && group.Status != types.RequestAccessGroupStatusPENDINGAPPROVAL {
		summary = fmt.Sprintf("%s %s %s's request", o.RequestReviewer.Email, statusLower, group.RequestedBy.Email)
	} else {
		summary = fmt.Sprintf("New request for %s from %s", group.AccessRuleSnapshot.Name, group.RequestedBy.Email)
	}

	start, _ := group.GetInterval(access.WithNow(clock.New().Now()))
	when := start.Format(time.Kitchen)

	requestDetails := []*slack.TextBlockObject{
		{
			Type: "mrkdwn",
			Text: fmt.Sprintf("*When:*\n%s", when),
		},
		{
			Type: "mrkdwn",
			Text: fmt.Sprintf("*Duration:*\n%s", group.RequestedTiming.Duration),
		},
		{
			Type: "mrkdwn",
			Text: fmt.Sprintf("*Status:*\n%s", statusNoUnder),
		},
	}

	// for _, v := range o.RequestArguments {
	// 	requestDetails = append(requestDetails, &slack.TextBlockObject{
	// 		Type: "mrkdwn",
	// 		Text: fmt.Sprintf("*%s:*\n%s", v.Title, v.Label),
	// 	})
	// }

	// Only show the Request reason if it is not empty
	if group.RequestPurposeReason != "" {
		requestDetails = append(requestDetails, &slack.TextBlockObject{
			Type: "mrkdwn",
			Text: fmt.Sprintf("*Request Reason:*\n%s", group.RequestPurposeReason),
		})
	}

	var richTextSummary string

	if o.IsWebhook && o.WasReviewed && group.Status != types.RequestAccessGroupStatusPENDINGAPPROVAL {
		richTextSummary = fmt.Sprintf("*%s %s %s's request*", o.RequestorEmail, statusLower, group.RequestedBy.Email)
	} else {
		richTextSummary = fmt.Sprintf("*<%s|New request for %s> from %s*", o.ReviewURLs.Review, group.AccessRuleSnapshot.Name, requestor)
	}

	msg = slack.NewBlockMessage(
		slack.SectionBlock{
			Type: slack.MBTSection,
			Text: &slack.TextBlockObject{
				Type: slack.MarkdownType,
				Text: richTextSummary,
			},
		},
		slack.SectionBlock{
			Type:   slack.MBTSection,
			Fields: requestDetails,
		},
	)

	if o.WasReviewed || group.Status == types.RequestAccessGroupStatusDECLINED {
		t := time.Now()
		when := types.ExpiryString(t)

		text := fmt.Sprintf("*Reviewed by* %s at %s", o.RequestorEmail, when)

		if group.Status == types.RequestAccessGroupStatusDECLINED {
			text = fmt.Sprintf("*Declined by* %s at %s", o.RequestReviewer.Email, when)
		}

		reviewContextBlock := slack.NewContextBlock("", slack.TextBlockObject{
			Type: slack.MarkdownType,
			Text: text,
		})

		msg.Blocks.BlockSet = append(msg.Blocks.BlockSet, reviewContextBlock)
	}

	// If the request has just been sent (PENDING), then append Action Blocks
	if group.Status == types.RequestAccessGroupStatusPENDINGAPPROVAL && !isCancelledOrRevoked {
		msg.Blocks.BlockSet = append(msg.Blocks.BlockSet, slack.NewActionBlock("review_actions",
			slack.ButtonBlockElement{
				Type:     slack.METButton,
				Text:     &slack.TextBlockObject{Type: slack.PlainTextType, Text: "Approve"},
				Style:    slack.StylePrimary,
				ActionID: "approve",
				Value:    "approve",
				URL:      o.ReviewURLs.Approve,
			},
			slack.ButtonBlockElement{
				Type:     slack.METButton,
				Text:     &slack.TextBlockObject{Type: slack.PlainTextType, Text: "Close Request"},
				Style:    slack.StyleDanger,
				ActionID: "deny",
				Value:    "deny",
				URL:      o.ReviewURLs.Deny,
			},
		))

	}

	return summary, msg
}

type RequestDetailMessageOpts struct {
	Request access.GroupWithTargets
	// the message that renders in the header of the slack message
	HeadingMessage string
}

// Builds a contextual request detail message, with an optional HeadingMessage to be rendered in the header, this is fired after a request has been reviewed
func BuildRequestDetailMessage(o RequestDetailMessageOpts) (summary string, msg slack.Message) {

	accessGroup := o.Request.Group

	summary = fmt.Sprintf("Request detail for %s", accessGroup.AccessRuleSnapshot.Name)

	// when := types.ExpiryString(expires)
	start, end := accessGroup.GetInterval(access.WithNow(clock.New().Now()))

	duration := end.Sub(start) // if this is off by a couple seconds it could make the duration values inconsistent

	requestDetails := []*slack.TextBlockObject{
		{
			Type: "mrkdwn",
			Text: fmt.Sprintf("*Duration:*\n%s", duration),
		},
		{
			Type: "mrkdwn",
			Text: fmt.Sprintf("*Expires:*\n%s", end),
		},
	}

	// FIELD MAPPINGS
	//
	// If we want to map over fields we can use `o.Request.Targets`
	//
	// 🚨🚨 Since a request group contains multipe targets, each with multiple fields 🚨🚨
	// How best should we dislay this in a single slack request
	// Is it ok to base it off the first target?
	// Is it ok to display a max of 2 fields from the target?

	// for _, v := range o.Request.Targets {
	// 	requestDetails = append(requestDetails, &slack.TextBlockObject{
	// 		Type: "mrkdwn",
	// 		Text: fmt.Sprintf("*%s:*\n%s", v.Title, v.Label),
	// 	})
	// }

	// Only show the Request reason if it is not empty
	if accessGroup.RequestPurposeReason != "" {
		requestDetails = append(requestDetails, &slack.TextBlockObject{
			Type: "mrkdwn",
			Text: fmt.Sprintf("*Request Reason:*\n%s", accessGroup.RequestPurposeReason),
		})
	}

	msg = slack.NewBlockMessage(
		slack.SectionBlock{
			Type: slack.MBTSection,
			Text: &slack.TextBlockObject{
				Type: slack.MarkdownType,
				Text: o.HeadingMessage,
			},
		},
		slack.SectionBlock{
			Type:   slack.MBTSection,
			Fields: requestDetails,
		},
	)

	return summary, msg
}
