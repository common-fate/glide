package slacknotifier

import (
	"fmt"
	"strings"
	"time"

	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/common-fate/common-fate/pkg/notifiers"
	"github.com/common-fate/common-fate/pkg/rule"
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
	Request          access.Request
	RequestArguments []types.With
	Rule             rule.AccessRule
	ReviewURLs       notifiers.ReviewURLs
	// optional field that will replace the default requestor email with a slack @mention
	RequestorSlackID string
	RequestorEmail   string
	WasReviewed      bool
	RequestReviewer  *identity.User
	//Optional field for a user or group to be tagged in the message.
	TaggedUser string
	IsWebhook  bool
}

/**
 * BuildRequestReviewMessage builds a slack message for a request review based on the contextual RequestMessageOpts
 */
func BuildRequestReviewMessage(o RequestMessageOpts) (summary string, msg slack.Message) {
	requestor := o.RequestorEmail
	if o.RequestorSlackID != "" {
		requestor = fmt.Sprintf("<@%s>", o.RequestorSlackID)
	}

	status := titleCase(string(o.Request.Status))
	statusLower := strings.ToLower(status)

	if o.IsWebhook && o.WasReviewed && o.Request.Status != access.PENDING {
		summary = fmt.Sprintf("%s %s %s's request", o.RequestReviewer.Email, statusLower, o.RequestorEmail)
	} else {
		summary = fmt.Sprintf("New request for %s from %s", o.Rule.Name, o.RequestorEmail)
	}

	when := "ASAP"
	if o.Request.RequestedTiming.StartTime != nil {
		t := o.Request.RequestedTiming.StartTime
		when = types.ExpiryString(*t)
	}

	requestDetails := []*slack.TextBlockObject{
		{
			Type: "mrkdwn",
			Text: fmt.Sprintf("*When:*\n%s", when),
		},
		{
			Type: "mrkdwn",
			Text: fmt.Sprintf("*Duration:*\n%s", o.Request.RequestedTiming.Duration),
		},
		{
			Type: "mrkdwn",
			Text: fmt.Sprintf("*Status:*\n%s", status),
		},
	}

	for _, v := range o.RequestArguments {
		requestDetails = append(requestDetails, &slack.TextBlockObject{
			Type: "mrkdwn",
			Text: fmt.Sprintf("*%s:*\n%s", v.Title, v.Label),
		})
	}

	// Only show the Request reason if it is not empty
	if o.Request.Data.Reason != nil && len(*o.Request.Data.Reason) > 0 {
		requestDetails = append(requestDetails, &slack.TextBlockObject{
			Type: "mrkdwn",
			Text: fmt.Sprintf("*Request Reason:*\n%s", *o.Request.Data.Reason),
		})
	}

	//If a tagged user is specified then add it to the message.
	if o.TaggedUser != "" {
		requestDetails = append(requestDetails, &slack.TextBlockObject{
			Type: "mrkdwn",
			Text: fmt.Sprintf("*Approver:*\n%s", o.TaggedUser),
		})
	}

	var richTextSummary string

	if o.IsWebhook && o.WasReviewed && o.Request.Status != access.PENDING {
		richTextSummary = fmt.Sprintf("*%s %s %s's request*", o.RequestReviewer.Email, statusLower, o.RequestorEmail)
	} else {
		richTextSummary = fmt.Sprintf("*<%s|New request for %s> from %s*", o.ReviewURLs.Review, o.Rule.Name, requestor)
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

	if o.WasReviewed || o.Request.Status == access.CANCELLED {
		t := time.Now()
		when := types.ExpiryString(t)

		text := fmt.Sprintf("*Reviewed by* %s at %s", o.RequestReviewer.Email, when)

		if o.Request.Status == access.CANCELLED {
			text = fmt.Sprintf("*Cancelled by* %s at %s", o.RequestorEmail, when)
		}

		reviewContextBlock := slack.NewContextBlock("", slack.TextBlockObject{
			Type: slack.MarkdownType,
			Text: text,
		})

		msg.Blocks.BlockSet = append(msg.Blocks.BlockSet, reviewContextBlock)
	}

	// If the request has just been sent (PENDING), then append Action Blocks
	if o.Request.Status == access.PENDING {
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
	Request          access.Request
	RequestArguments []types.With
	Rule             rule.AccessRule
	// the message that renders in the header of the slack message
	HeadingMessage string
}

// Builds a contextual request detail message, with an optional HeadingMessage to be rendered in the header, this is fired after a request has been reviewed
func BuildRequestDetailMessage(o RequestDetailMessageOpts) (summary string, msg slack.Message) {

	summary = fmt.Sprintf("Request detail for %s", o.Rule.Name)

	var expires time.Time
	if o.Request.RequestedTiming.StartTime != nil {
		expires = o.Request.RequestedTiming.StartTime.Add(o.Request.RequestedTiming.Duration)
	} else {
		expires = time.Now().Add(o.Request.RequestedTiming.Duration)
	}

	when := types.ExpiryString(expires)

	requestDetails := []*slack.TextBlockObject{

		{
			Type: "mrkdwn",
			Text: fmt.Sprintf("*Duration:*\n%s", o.Request.RequestedTiming.Duration),
		},
		{
			Type: "mrkdwn",
			Text: fmt.Sprintf("*Expires:*\n%s", when),
		},
	}

	for _, v := range o.RequestArguments {
		requestDetails = append(requestDetails, &slack.TextBlockObject{
			Type: "mrkdwn",
			Text: fmt.Sprintf("*%s:*\n%s", v.Title, v.Label),
		})
	}

	// Only show the Request reason if it is not empty
	if o.Request.Data.Reason != nil && len(*o.Request.Data.Reason) > 0 {
		requestDetails = append(requestDetails, &slack.TextBlockObject{
			Type: "mrkdwn",
			Text: fmt.Sprintf("*Request Reason:*\n%s", *o.Request.Data.Reason),
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
