package slacknotifier

import (
	"fmt"
	"strings"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/common-fate/common-fate/pkg/access"
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
	Group         access.Group
	RequestReason string
	// AccessGroups []access.Group
	// RequestArguments []types.With

	ReviewURLs notifiers.ReviewURLs
	// optional field that will replace the default requestor email with a slack @mention
	RequestorSlackID string
	RequestorEmail   string
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

	// @TODO: we're wondering where the best place it to run this itteration logic...

	requestor := group.RequestedBy.Email
	if o.RequestorSlackID != "" {
		requestor = fmt.Sprintf("<@%s>", o.RequestorSlackID)
	}

	status := titleCase(string(group.Status))
	statusLower := strings.ToLower(status)

	if o.IsWebhook && o.WasReviewed && group.Status != types.RequestAccessGroupStatusPENDINGAPPROVAL {
		summary = fmt.Sprintf("%s %s %s's request", group, statusLower, group.RequestedBy.Email)
	} else {
		summary = fmt.Sprintf("New request for %s from %s", group.AccessRuleSnapshot.Name, group.RequestedBy.Email)
	}

	requestDetails := []*slack.TextBlockObject{
		{
			Type: "mrkdwn",
			Text: fmt.Sprintf("*When:*\n%s", group.RequestedTiming.StartTime),
		},
		{
			Type: "mrkdwn",
			Text: fmt.Sprintf("*Duration:*\n%s", group.RequestedTiming.Duration),
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

	if o.RequestReason != "" {
		requestDetails = append(requestDetails, &slack.TextBlockObject{
			Type: "mrkdwn",
			Text: fmt.Sprintf("*Request Reason:*\n%s", o.RequestReason),
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
			text = fmt.Sprintf("*Cancelled by* %s at %s", group.RequestedBy.Email, when)
		}

		reviewContextBlock := slack.NewContextBlock("", slack.TextBlockObject{
			Type: slack.MarkdownType,
			Text: text,
		})

		msg.Blocks.BlockSet = append(msg.Blocks.BlockSet, reviewContextBlock)
	}

	// If the request has just been sent (PENDING), then append Action Blocks
	if group.Status == types.RequestAccessGroupStatusPENDINGAPPROVAL {
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

// func (n *SlackNotifier) SendUpdatesForAccess(ctx context.Context, log *zap.SugaredLogger, request access.Request) {
// 	// Loop over the request reviewers
// 	reviewers := storage.ListRequestReviewers{RequestID: request.ID}
// 	_, err := n.DB.Query(ctx, &reviewers)
// 	if err != nil && err != ddb.ErrNoItems {
// 		log.Errorw("failed to fetch reviewers for request", zap.Error(err))
// 		return
// 	}
// 	reqReviewer := storage.GetUser{ID: requestEvent.ReviewerID}
// 	_, err = n.DB.Query(ctx, &reqReviewer)
// 	if err != nil && request.Status != access.CANCELLED {
// 		log.Errorw("failed to fetch reviewer for request which wasn't cancelled", zap.Error(err))
// 		return
// 	}
// 	reviewURL, err := notifiers.ReviewURL(n.FrontendURL, request.ID)
// 	if err != nil {
// 		log.Errorw("building review URL", zap.Error(err))
// 		return
// 	}
// 	requestArguments, err := n.RenderRequestArguments(ctx, log, request, rule)
// 	if err != nil {
// 		log.Errorw("failed to generate request arguments, skipping including them in the slack message", "error", err)
// 	}
// 	log.Infow("messaging reviewers", "reviewers", reviewers.Result)
// 	if n.directMessageClient != nil {
// 		// get the requestor's Slack user ID if it exists to render it nicely in the message to approvers.
// 		var slackUserID string
// 		requestor, err := n.directMessageClient.client.GetUserByEmailContext(ctx, requestingUser.Email)
// 		if err != nil {
// 			// log this instead of returning
// 			log.Errorw("failed to get slack user id, defaulting to email", "user", requestingUser.Email, zap.Error(err))
// 		}
// 		if requestor != nil {
// 			slackUserID = requestor.ID
// 		}
// 		_, msg := BuildRequestReviewMessage(RequestMessageOpts{
// 			Request:          request,
// 			RequestArguments: requestArguments,
// 			Rule:             rule,
// 			RequestorSlackID: slackUserID,
// 			RequestorEmail:   requestingUser.Email,
// 			ReviewURLs:       reviewURL,
// 			WasReviewed:      true,
// 			RequestReviewer:  reqReviewer.Result,
// 			IsWebhook:        false,
// 		})
// 		for _, usr := range reviewers.Result {
// 			err = n.UpdateMessageBlockForReviewer(ctx, usr, msg)
// 			if err != nil {
// 				log.Errorw("failed to update slack message", "user", usr, zap.Error(err))
// 			}
// 		}
// 	}

// 	// log for testing purposes
// 	if len(n.webhooks) > 0 {
// 		log.Infow("webhooks found", "webhooks", n.webhooks)
// 	}
// 	// this does not include the slackUserID because we don't have access to look it up
// 	summary, msg := BuildRequestReviewMessage(RequestMessageOpts{
// 		Request:          request,
// 		RequestArguments: requestArguments,
// 		Rule:             rule,
// 		RequestorEmail:   requestingUser.Email,
// 		ReviewURLs:       reviewURL,
// 		WasReviewed:      true,
// 		RequestReviewer:  reqReviewer.Result,
// 		IsWebhook:        true,
// 	})
// 	for _, webhook := range n.webhooks {
// 		err = webhook.SendWebhookMessage(ctx, msg.Blocks, summary)
// 		if err != nil {
// 			log.Errorw("failed to send review message to incomingWebhook channel", "error", err)
// 		}
// 	}
// }
