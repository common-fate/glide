package slacknotifier

import (
	"fmt"
	"strings"
	"time"

	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/common-fate/common-fate/pkg/notifiers"
	"github.com/common-fate/common-fate/pkg/requests"
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
	Request          requests.Requestv2
	AccessGroups     []requests.AccessGroup
	RequestArguments []types.With

	ReviewURLs notifiers.ReviewURLs
	// optional field that will replace the default requestor email with a slack @mention
	RequestorSlackID string
	RequestorEmail   string
	WasReviewed      bool
	RequestReviewer  *identity.User
	IsWebhook        bool
}

/**
 * BuildRequestReviewMessage builds a slack message for a request review based on the contextual RequestMessageOpts
 */
func BuildRequestReviewMessage(o RequestMessageOpts) (summary string, msg slack.Message) {

	for _, access_group := range o.AccessGroups {
		requestor := o.RequestorEmail
		if o.RequestorSlackID != "" {
			requestor = fmt.Sprintf("<@%s>", o.RequestorSlackID)
		}

		status := titleCase(string(access_group.Status))
		statusLower := strings.ToLower(status)

		if o.IsWebhook && o.WasReviewed && access_group.Status != requests.PENDING {
			summary = fmt.Sprintf("%s %s %s's request", o.RequestReviewer.Email, statusLower, o.RequestorEmail)
		} else {
			summary = fmt.Sprintf("New request for %s from %s", access_group.AccessRule.Name, o.RequestorEmail)
		}

		when := "ASAP"
		if access_group.TimeConstraints.StartTime != nil {
			t := access_group.TimeConstraints.StartTime
			when = types.ExpiryString(*t)
		}

		requestDetails := []*slack.TextBlockObject{
			{
				Type: "mrkdwn",
				Text: fmt.Sprintf("*When:*\n%s", when),
			},
			{
				Type: "mrkdwn",
				Text: fmt.Sprintf("*Duration:*\n%s", access_group.TimeConstraints.Duration),
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
		if o.Request.Context.Reason != nil && len(*o.Request.Context.Reason) > 0 {
			requestDetails = append(requestDetails, &slack.TextBlockObject{
				Type: "mrkdwn",
				Text: fmt.Sprintf("*Request Reason:*\n%s", *o.Request.Context.Reason),
			})
		}

		var richTextSummary string

		if o.IsWebhook && o.WasReviewed && access_group.Status != requests.PENDING {
			richTextSummary = fmt.Sprintf("*%s %s %s's request*", o.RequestReviewer.Email, statusLower, o.RequestorEmail)
		} else {
			richTextSummary = fmt.Sprintf("*<%s|New request for %s> from %s*", o.ReviewURLs.Review, access_group.AccessRule.Name, requestor)
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

		if o.WasReviewed || access_group.Status == requests.CANCELLED {
			t := time.Now()
			when := types.ExpiryString(t)

			text := fmt.Sprintf("*Reviewed by* %s at %s", o.RequestReviewer.Email, when)

			if access_group.Status == requests.CANCELLED {
				text = fmt.Sprintf("*Cancelled by* %s at %s", o.RequestorEmail, when)
			}

			reviewContextBlock := slack.NewContextBlock("", slack.TextBlockObject{
				Type: slack.MarkdownType,
				Text: text,
			})

			msg.Blocks.BlockSet = append(msg.Blocks.BlockSet, reviewContextBlock)
		}

		// If the request has just been sent (PENDING), then append Action Blocks
		if access_group.Status == requests.PENDING {
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
	}

	return summary, msg
}

type RequestDetailMessageOpts struct {
	AccessGroups []requests.AccessGroup
	Request      requests.Requestv2
	// the message that renders in the header of the slack message
	HeadingMessage string
}

// Builds a contextual request detail message, with an optional HeadingMessage to be rendered in the header, this is fired after a request has been reviewed
func BuildRequestDetailMessage(o RequestDetailMessageOpts) (summary string, msg slack.Message) {

	for _, access_group := range o.AccessGroups {
		summary = fmt.Sprintf("Request detail for %s", access_group.AccessRule.Name)

		var expires time.Time
		if access_group.TimeConstraints.StartTime != nil {
			expires = access_group.TimeConstraints.StartTime.Add(access_group.TimeConstraints.Duration)
		} else {
			expires = time.Now().Add(access_group.TimeConstraints.Duration)
		}

		when := types.ExpiryString(expires)

		requestDetails := []*slack.TextBlockObject{

			{
				Type: "mrkdwn",
				Text: fmt.Sprintf("*Duration:*\n%s", access_group.TimeConstraints.Duration),
			},
			{
				Type: "mrkdwn",
				Text: fmt.Sprintf("*Expires:*\n%s", when),
			},
		}

		// for _, v := range access_group. {
		// 	requestDetails = append(requestDetails, &slack.TextBlockObject{
		// 		Type: "mrkdwn",
		// 		Text: fmt.Sprintf("*%s:*\n%s", v.Title, v.Label),
		// 	})
		// }

		// Only show the Request reason if it is not empty
		if o.Request.Context.Reason != nil && len(*o.Request.Context.Reason) > 0 {
			requestDetails = append(requestDetails, &slack.TextBlockObject{
				Type: "mrkdwn",
				Text: fmt.Sprintf("*Request Reason:*\n%s", *o.Request.Context.Reason),
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
	}

	return summary, msg
}
