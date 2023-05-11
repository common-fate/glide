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
	Request access.RequestWithGroupsWithTargets
	// AccessGroups []access.Group
	// RequestArguments []types.With

	ReviewURLs notifiers.ReviewURLs
	// optional field that will replace the default requestor email with a slack @mention
	RequestorSlackID string
	RequestorEmail   string
	WasReviewed      bool
	RequestReviewer  *identity.User
	IsWebhook        bool
}

/*
*
BuildRequestReviewMessage builds a slack message for a request review based on the contextual RequestMessageOpts
*/
func BuildRequestReviewMessage(o RequestMessageOpts) (summary string, msg slack.Message) {

	req := o.Request

	for _, access_group := range req.Groups {

		requestor := req.Request.RequestedBy.Email
		if o.RequestorSlackID != "" {
			requestor = fmt.Sprintf("<@%s>", o.RequestorSlackID)
		}

		status := titleCase(string(access_group.Group.Status))
		statusLower := strings.ToLower(status)

		if o.IsWebhook && o.WasReviewed && access_group.Group.Status != types.RequestAccessGroupStatusPENDINGAPPROVAL {
			summary = fmt.Sprintf("%s %s %s's request", o.RequestReviewer.Email, statusLower, req.Request.RequestedBy.Email)
		} else {
			summary = fmt.Sprintf("New request for %s from %s", access_group.Group.AccessRuleSnapshot.Name, req.Request.RequestedBy.Email)
		}

		// This can be deprecated bc ASAP now == no reviewers (and this is a reviewer message)
		// when := "ASAP"
		// 🚨🚨🚨 How is ASAP determined post-refactor? 🚨🚨🚨
		// if access_group.TimeConstraints.StartTime != nil {
		// 	t := access_group.TimeConstraints.StartTime
		// 	when = types.ExpiryString(*t)
		// }

		requestDetails := []*slack.TextBlockObject{
			{
				Type: "mrkdwn",
				Text: fmt.Sprintf("*When:*\n%s", access_group.Group.RequestedTiming.StartTime),
			},
			{
				Type: "mrkdwn",
				Text: fmt.Sprintf("*Duration:*\n%s", access_group.Group.RequestedTiming.Duration),
			},
			{
				Type: "mrkdwn",
				Text: fmt.Sprintf("*Status:*\n%s", status),
			},
		}

		// for _, v := range o.RequestArguments {
		// 	requestDetails = append(requestDetails, &slack.TextBlockObject{
		// 		Type: "mrkdwn",
		// 		Text: fmt.Sprintf("*%s:*\n%s", v.Title, v.Label),
		// 	})
		// }

		// Only show the Request reason if it is not empty

		if req.Request.Purpose.Reason != nil && len(*req.Request.Purpose.Reason) > 0 {
			requestDetails = append(requestDetails, &slack.TextBlockObject{
				Type: "mrkdwn",
				Text: fmt.Sprintf("*Request Reason:*\n%s", *req.Request.Purpose.Reason),
			})
		}

		var richTextSummary string

		if o.IsWebhook && o.WasReviewed && access_group.Group.Status != types.RequestAccessGroupStatusPENDINGAPPROVAL {
			richTextSummary = fmt.Sprintf("*%s %s %s's request*", o.RequestReviewer.Email, statusLower, req.Request.RequestedBy.Email)
		} else {
			richTextSummary = fmt.Sprintf("*<%s|New request for %s> from %s*", o.ReviewURLs.Review, access_group.Group.AccessRuleSnapshot.Name, requestor)
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

		if o.WasReviewed || access_group.Group.Status == types.RequestAccessGroupStatusDECLINED {
			t := time.Now()
			when := types.ExpiryString(t)

			text := fmt.Sprintf("*Reviewed by* %s at %s", o.RequestReviewer.Email, when)

			if access_group.Group.Status == types.RequestAccessGroupStatusDECLINED {
				text = fmt.Sprintf("*Cancelled by* %s at %s", req.Request.RequestedBy.Email, when)
			}

			reviewContextBlock := slack.NewContextBlock("", slack.TextBlockObject{
				Type: slack.MarkdownType,
				Text: text,
			})

			msg.Blocks.BlockSet = append(msg.Blocks.BlockSet, reviewContextBlock)
		}

		// If the request has just been sent (PENDING), then append Action Blocks
		if access_group.Group.Status == types.RequestAccessGroupStatusPENDINGAPPROVAL {
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
	Request access.RequestWithGroupsWithTargets
	// the message that renders in the header of the slack message
	HeadingMessage string
}

// Builds a contextual request detail message, with an optional HeadingMessage to be rendered in the header, this is fired after a request has been reviewed
func BuildRequestDetailMessage(o RequestDetailMessageOpts) (summary string, msg slack.Message) {

	req := o.Request

	for _, access_group := range req.Groups {
		summary = fmt.Sprintf("Request detail for %s", access_group.Group.AccessRuleSnapshot.Name)
		/**
		var expires time.Time
		// has start time...
		if access_group.Group.RequestedTiming.StartTime != nil {
			// has override..
			if access_group.Group.OverrideTiming.StartTime != nil {
				// has override duration
				if access_group.Group.OverrideTiming.Duration != nil {
					expires = access_group.Group.OverrideTiming.StartTime.Add(*access_group.Group.OverrideTiming.Duration)
				} else {
					expires = access_group.Group.OverrideTiming.StartTime.Add(access_group.Group.OverrideTiming.Duration)
				}
			} else {
				// has override duration
				expires = access_group.Group.RequestedTiming.StartTime.Add(access_group.Group.RequestedTiming.Duration)
			}
		} else {
			expires = time.Now().Add(access_group.TimeConstraints.Duration)
		}


		...

		We had this but we probably want to leveraage FinalTiming if possibel to reduce if/else logic

		*/

		// has start time...
		// if access_group.Group.RequestedTiming.StartTime != nil {
		// 	// has override..
		// 	if access_group.Group.OverrideTiming.StartTime != nil {
		// 		// has override duration
		// 		if access_group.Group.OverrideTiming.Duration != nil {
		// 			expires = access_group.Group.OverrideTiming.StartTime.Add(*access_group.Group.OverrideTiming.Duration)
		// 		} else {
		// 			expires = access_group.Group.OverrideTiming.StartTime.Add(access_group.Group.OverrideTiming.Duration)
		// 		}
		// 	} else {
		// 		// has override duration
		// 		expires = access_group.Group.RequestedTiming.StartTime.Add(access_group.Group.RequestedTiming.Duration)
		// 	}
		// } else {
		// 	expires = time.Now().Add(access_group.Group.RequestedTiming.Duration)
		// }

		// var expires time.Time

		// if access_group.Group.OverrideTiming != nil {
		// 	expires = access_group.Group.OverrideTiming.StartTime.Add(access_group.Group.OverrideTiming.Duration)
		// } else {
		// 	expires = access_group.Group.RequestedTiming.StartTime.Add(access_group.Group.RequestedTiming.Duration)
		// }

		// clock := clock.New().Now()

		// when := types.ExpiryString(expires)
		start, end := access_group.Group.GetInterval(access.WithNow(clock.New().Now()))

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

		// for _, v := range access_group. {
		// 	requestDetails = append(requestDetails, &slack.TextBlockObject{
		// 		Type: "mrkdwn",
		// 		Text: fmt.Sprintf("*%s:*\n%s", v.Title, v.Label),
		// 	})
		// }

		// Only show the Request reason if it is not empty
		if req.Request.Purpose.Reason != nil && len(*req.Request.Purpose.Reason) > 0 {
			requestDetails = append(requestDetails, &slack.TextBlockObject{
				Type: "mrkdwn",
				Text: fmt.Sprintf("*Request Reason:*\n%s", *req.Request.Purpose.Reason),
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
