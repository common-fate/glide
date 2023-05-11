package slacknotifier

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/gevent"
	"github.com/common-fate/common-fate/pkg/storage"
	"go.uber.org/zap"
)

func (n *SlackNotifier) HandleAccessGroupEvent(ctx context.Context, log *zap.SugaredLogger, event events.CloudWatchEvent) error {

	switch event.DetailType {
	case gevent.AccessGroupApprovedType:

		var accessGroupEvent gevent.AccessGroupApproved
		err := json.Unmarshal(event.Detail, &accessGroupEvent)
		if err != nil {
			return err
		}
		accessGroup := accessGroupEvent.AccessGroup

		// accessGroup.Group == `access.Group`
		// can fetched request using an id

		// REQUESTOR Message:
		// "your access to X no. of resources for Y access rule has been approved"
		msg := fmt.Sprintf(":white_check_mark: Your request to access *%s* has been approved.", accessGroup.Group.AccessRuleSnapshot.Name)
		fallback := fmt.Sprintf("Your request to access %s has been approved.", accessGroup.Group.AccessRuleSnapshot.Name)
		n.sendAccessGroupDetailsMessage(ctx, log, accessGroup, msg, fallback)

		// REVIWER Message Update:
		// n.SendUpdatesForRequest(ctx, log, request, requestEvent, requestedRule, requestingUserQuery.Result)

	case gevent.AccessGroupDeclinedType:

		var accessGroupEvent gevent.AccessGroupApproved
		err := json.Unmarshal(event.Detail, &accessGroupEvent)
		if err != nil {
			return err
		}
		accessGroup := accessGroupEvent.AccessGroup

		// REQUESTOR Message:
		// "your access to X no. of resources for Y access rule has been declined"
		msg := fmt.Sprintf("Your request to access *%s* has been declined.", accessGroup.Group.AccessRuleSnapshot.Name)
		fallback := fmt.Sprintf("Your request to access %s has been declined.", accessGroup.Group.AccessRuleSnapshot.Name)
		n.sendAccessGroupDetailsMessage(ctx, log, accessGroup, msg, fallback)

		// REVIWER Message Update:
		n.sendAccessGroupUpdates(ctx, log, accessGroup)

	default:
		zap.S().Infow("unhandled access group event", "detailType", event.DetailType)
	}

	return nil
}

// sendAccessGroupDetailsMessage sends a message to the user who requested access with details about the request. Sent only on access create/approved
func (n *SlackNotifier) sendAccessGroupDetailsMessage(ctx context.Context, log *zap.SugaredLogger, accessGroup access.GroupWithTargets, headingMsg string, summary string) {

	approvalRequired := false // TODO

	var HAS_SLACK_CLIENT = n.directMessageClient != nil
	var HAS_SLACK_WEBHOOKS = len(n.webhooks) > 0

	requestor := accessGroup.Group.RequestedBy
	// reviewers := accessGroup.Group.GroupReviewers

	// 🚨🚨 `requestedRule.Name` references to -> a count of the Resource Gropus that the user is requesting access to

	if HAS_SLACK_CLIENT {
		_, msg := BuildRequestDetailMessage(RequestDetailMessageOpts{
			Request:        accessGroup,
			HeadingMessage: headingMsg,
		})

		_, err := SendMessageBlocks(ctx, n.directMessageClient.client, requestor.Email, msg, summary)

		if err != nil {
			log.Errorw("failed to send slack message", "user", requestor, zap.Error(err))
		}
	}

	if HAS_SLACK_WEBHOOKS {
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

/*
What we need to do:
- Either adapt `BuildRequestDetailMessage` to dual serve both `access.Request` and `access.GroupWithTargets` or create a new function for each
- Could also do with a major cleanup commit of the code base

See if you can reproduce BuildRequestDetailMessage with with just those props



*/

func (n *SlackNotifier) sendAccessGroupUpdates(ctx context.Context, log *zap.SugaredLogger, accessGroup access.GroupWithTargets) {

	var HAS_SLACK_CLIENT = n.directMessageClient != nil
	var HAS_SLACK_WEBHOOKS = len(n.webhooks) > 0

	requestor := accessGroup.Group.RequestedBy

	if HAS_SLACK_CLIENT {

		for _, reviewer := range accessGroup.Group.RequestReviewers {

			reqReviewer := storage.GetRequestReviewer{
				RequestID:  accessGroup.Group.RequestID,
				ReviewerID: reviewer,
			}
			_, err := n.DB.Query(ctx, &reqReviewer)
			if err != nil {
				log.Errorw("failed to get request reviewer", "error", err)
				return
			}

			// reqReviewer.Result.
			// storage.GetRequestGroupTarget

			// TODO: wondering if we can access request here since it will be needed to construct review messages
			// accessGroup.Group.RequestedBy

			// summ, slackMsg := BuildRequestReviewMessage(RequestMessageOpts{
			// 	Group: accessGroup.Group,
			// })

			// 🚨🚨🚨🚨 TODO: now fire this off 🚨🚨🚨🚨

			summary, slackMsg := BuildRequestReviewMessage(RequestMessageOpts{
				Group: accessGroup.Group,
			})

			_, err = SendMessageBlocks(ctx, n.directMessageClient.client, requestor.Email, slackMsg, summary)

			if err != nil {
				log.Errorw("failed to send slack message", "user", requestor, zap.Error(err))
			}
		}

	}

	if HAS_SLACK_WEBHOOKS {
		// 	for _, webhook := range n.webhooks {
		// 		// if !approvalRequired {
		// 		// headingMsg = fmt.Sprintf(":white_check_mark: %s's request to access *%s* has been automatically approved.\n", requestingUser.Email, requestedRule.Name)

		// 		// summary = fmt.Sprintf("%s's request to access %s has been automatically approved.", requestingUser.Email, requestedRule.Name)
		// 		// }
		// 		_, msg := BuildRequestDetailMessage(RequestDetailMessageOpts{
		// 			Request: accessGroup,
		// 			// RequestArguments: requestArguments,
		// 			HeadingMessage: headingMsg,
		// 		})

		// 		err := webhook.SendWebhookMessage(ctx, msg.Blocks, summary)
		// 		if err != nil {
		// 			log.Errorw("failed to send slack message to webhook channel", "error", err)
		// 		}
		// 	}
	}

}
