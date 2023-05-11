package slacknotifier

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/gevent"
	"go.uber.org/zap"
)

// 🚨🚨🚨🚨🚨🚨🚨🚨🚨🚨🚨🚨 TODO: replace me to adopt accessGroup event type 🚨🚨🚨🚨🚨🚨🚨🚨🚨🚨🚨🚨🚨🚨🚨🚨🚨🚨

func (n *SlackNotifier) HandleAccessGroupEvent(ctx context.Context, log *zap.SugaredLogger, event events.CloudWatchEvent) error {

	var HAS_SLACK_CLIENT = n.directMessageClient != nil
	var HAS_SLACK_WEBHOOKS = len(n.webhooks) > 0

	// get the message text based on the event type
	switch event.DetailType {
	case gevent.AccessGroupApprovedType:

		var accessGroupEvent gevent.AccessGroupApproved
		err := json.Unmarshal(event.Detail, &accessGroupEvent)
		if err != nil {
			return err
		}
		accessGroup := accessGroupEvent.AccessGroup

		// Getting the requestor...
		// accessGroupEvent.AccessGroup.Group.RequestedBy

		// Getting the reviewers and their notifications to update....
		// accessGroupEvent.AccessGroup.Group.GroupReviewers.
		// storage.GetRequestReviewer{
		// 	ReviewerID: ,
		// }

		// REQUESTOR Message:
		// "your access to X no. of resources for Y access rule has been approved"

		msg := fmt.Sprintf(":white_check_mark: Your request to access *%s* has been approved.", accessGroup.Group.AccessRuleSnapshot.Name)
		fallback := fmt.Sprintf("Your request to access %s has been approved.", accessGroup.Group.AccessRuleSnapshot.Name)
		//
		// 🚨🚨 sendRequestDetailsMessage should be sent to the requestor only and it should go on at RequestCreated event
		n.sendRequestDetailsMessage(ctx, log, accessGroup, msg, fallback)
		// n.SendUpdatesForRequest(ctx, log, request, requestEvent, requestedRule, requestingUserQuery.Result)

		// REVIWER Message Update:

	case gevent.AccessGroupDeclinedType:
		zap.S().Infow("unhandled grant event", "detailType", event.DetailType)
	}
	// if msg != "" {
	// 	_ = fallback
	// 	_ = accessory
	// 	// _, err = SendMessage(ctx, n.directMessageClient.client, gq.Result.RequestedBy.Email, msg, fallback, accessory)
	// 	// return err
	// }
	return nil
}

// sendAccessGroupDetailsMessage sends a message to the user who requested access with details about the request. Sent only on access create/approved
func (n *SlackNotifier) sendAccessGroupDetailsMessage(ctx context.Context, log *zap.SugaredLogger, accessGroup access.GroupWithTargets, headingMsg string, summary string) {
	// requestArguments, err := n.RenderRequestArguments(ctx, log, request, requestedRule)
	// if err != nil {
	// 	log.Errorw("failed to generate request arguments, skipping including them in the slack message", "error", err)
	// }

	// var slackUserID string
	// slackRequestor, err := n.directMessageClient.client.GetUserByEmailContext(ctx, "jordi@commonfate.io")
	// requestor, err := n.directMessageClient.client.GetUserByEmailContext(ctx, accessGroup.Group.RequestedBy.Email)

	approvalRequired := false // TODO

	requestor := accessGroup.Group.RequestedBy

	// 🚨🚨 `requestedRule.Name` references to -> a count of the Resource Gropus that the user is requesting access to

	if n.directMessageClient != nil || len(n.webhooks) > 0 {
		if n.directMessageClient != nil {
			_, msg := BuildRequestDetailMessage(RequestDetailMessageOpts{
				Request:        accessGroup,
				HeadingMessage: headingMsg,
			})

			_, err := SendMessageBlocks(ctx, n.directMessageClient.client, requestor.Email, msg, summary)

			if err != nil {
				log.Errorw("failed to send slack message", "user", requestor, zap.Error(err))
			}
		}

		for _, webhook := range n.webhooks {
			if !approvalRequired {
				// headingMsg = fmt.Sprintf(":white_check_mark: %s's request to access *%s* has been automatically approved.\n", requestingUser.Email, requestedRule.Name)

				// summary = fmt.Sprintf("%s's request to access %s has been automatically approved.", requestingUser.Email, requestedRule.Name)
			}
			_, msg := BuildRequestDetailMessage(RequestDetailMessageOpts{
				Request: accessGroup,
				// RequestArguments: requestArguments,
				HeadingMessage: headingMsg,
			})

			err := webhook.SendWebhookMessage(ctx, msg.Blocks, summary)
			if err != nil {
				log.Errorw("failed to send slack message to webhook channel", "error", err)
			}
		}
	}
}
