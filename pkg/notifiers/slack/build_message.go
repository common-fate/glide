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
	RequestorSlackID string
	RequestorEmail   string
	WasReviewed      bool
	RequestReviewer  *identity.User
	IsWebhook        bool
}

func BuildRequestMessage(o RequestMessageOpts) (summary string, msg slack.Message) {
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
		when = fmt.Sprintf("<!date^%d^{date_short_pretty} at {time}|%s>", t.Unix(), t.String())
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
		when = fmt.Sprintf("<!date^%d^{date_short_pretty} at {time}|%s>", t.Unix(), t.String())

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
	RequestorSlackID string
	RequestorEmail   string
	IsWebhook        bool
	OriginalMessage  string
}

// Builds a generic rule detail table to be included in any message
func BuildRequestDetailMessage(o RequestDetailMessageOpts) (summary string, msg slack.Message) {

	status := titleCase(string(o.Request.Status))

	summary = fmt.Sprintf("Request detail for %s", o.Rule.Name)

	when := "ASAP"
	if o.Request.RequestedTiming.StartTime != nil {
		t := o.Request.RequestedTiming.StartTime
		when = fmt.Sprintf("<!date^%d^{date_short_pretty} at {time}|%s>", t.Unix(), t.String())
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

	msg = slack.NewBlockMessage(
		slack.SectionBlock{
			Type: slack.MBTSection,
			Text: &slack.TextBlockObject{
				Type: slack.MarkdownType,
				Text: o.OriginalMessage,
			},
		},
		slack.SectionBlock{
			Type:   slack.MBTSection,
			Fields: requestDetails,
		},
	)

	return summary, msg
}
