package slacknotifier

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/gevent"
	"github.com/common-fate/common-fate/pkg/notifiers"
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

		// REQUESTOR Message:
		// "your access to X no. of resources for Y access rule has been approved"
		msg := fmt.Sprintf(":white_check_mark: Your request to access *%s* has been approved.", accessGroup.Group.AccessRuleSnapshot.Name)
		fallback := fmt.Sprintf("Your request to access %s has been approved.", accessGroup.Group.AccessRuleSnapshot.Name)
		n.sendAccessGroupDetailsMessageRequestor(ctx, log, accessGroup, msg, fallback)

		// REVIEWER Message Update:
		n.sendAccessGroupUpdatesReviewer(ctx, log, accessGroup)

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
		n.sendAccessGroupDetailsMessageRequestor(ctx, log, accessGroup, msg, fallback)

		// REVIEWER Message Update:
		n.sendAccessGroupUpdatesReviewer(ctx, log, accessGroup)

	default:
		zap.S().Infow("unhandled access group event", "detailType", event.DetailType)
	}

	return nil
}

// sendAccessGroupDetailsMessageRequestor sends a message to the Requestor with details about the request. Sent only on AccessGroupDeclinedType, AccessGroupApprovedType
func (n *SlackNotifier) sendAccessGroupDetailsMessageRequestor(ctx context.Context, log *zap.SugaredLogger, accessGroup access.GroupWithTargets, headingMsg string, summary string) {

	var HAS_SLACK_CLIENT = n.directMessageClient != nil
	var HAS_SLACK_WEBHOOKS = len(n.webhooks) > 0

	// This is used for dev testing puprposes only,
	// this allows the requestor to act as a reviewer so you can test both,
	// notification types in one go.
	var OVERRIDE_DEV = false
	var OVERRIDE_EMAIL = "jordi@commonfate.io"

	requestor := accessGroup.Group.RequestedBy

	if HAS_SLACK_CLIENT {
		_, msg := BuildRequestDetailMessage(RequestDetailMessageOpts{
			Request:        accessGroup,
			HeadingMessage: headingMsg,
		})

		if OVERRIDE_DEV {
			requestor.Email = OVERRIDE_EMAIL
		}

		_, err := SendMessageBlocks(ctx, n.directMessageClient.client, requestor.Email, msg, summary)

		if err != nil {
			log.Errorw("failed to send slack message", "user", requestor, zap.Error(err))
		}
	}

	if HAS_SLACK_WEBHOOKS {
		for _, webhook := range n.webhooks {
			_, msg := BuildRequestDetailMessage(RequestDetailMessageOpts{
				Request:        accessGroup,
				HeadingMessage: headingMsg,
			})

			err := webhook.SendWebhookMessage(ctx, msg.Blocks, summary)
			if err != nil {
				log.Errorw("failed to send slack message to webhook channel", "error", err)
			}
		}
	}
}

func (n *SlackNotifier) sendAccessGroupUpdatesReviewer(ctx context.Context, log *zap.SugaredLogger, accessGroup access.GroupWithTargets) {

	var HAS_SLACK_CLIENT = n.directMessageClient != nil
	// var HAS_SLACK_WEBHOOKS = len(n.webhooks) > 0

	// This is used for dev testing puprposes only,
	// this allows the requestor to act as a reviewer so you can test both,
	// notification types in one go.
	var OVERRIDE_DEV = false
	var OVERRIDE_EMAIL = "jordi@commonfate.io"

	requestor := accessGroup.Group.RequestedBy

	if HAS_SLACK_CLIENT {

		// Loop over the request reviewers...
		for _, reviewer := range accessGroup.Group.RequestReviewers {

			reqReviewer := storage.GetRequestReviewer{
				RequestID:  accessGroup.Group.RequestID,
				ReviewerID: reviewer,
			}
			_, err := n.DB.Query(ctx, &reqReviewer)
			if err != nil {
				log.Errorw("failed to get request reviewer", "error", err)
				continue
			}

			reviewURL, err := notifiers.ReviewURL(n.FrontendURL, accessGroup.Group.RequestID)
			if err != nil {
				log.Errorw("building review URL", zap.Error(err))
				return
			}

			var slackUserID string
			slackRequestor, err := n.directMessageClient.client.GetUserByEmailContext(ctx, requestor.Email)
			if err != nil {
				zap.S().Infow("couldn't get slack user from requestor - falling back to email address", "requestor.id", requestor.Email, zap.Error(err))
			}
			if slackRequestor != nil {
				slackUserID = slackRequestor.ID
			}

			reviewerUserObj := storage.GetUser{ID: reviewer}
			_, err = n.DB.Query(ctx, &reviewerUserObj)
			if err != nil {
				log.Errorw("failed to get reviewer user", "error", err)
				continue
			}

			if OVERRIDE_DEV {
				reviewerUserObj.Result.Email = OVERRIDE_EMAIL
			}

			_, slackMsg := BuildRequestReviewMessage(RequestMessageOpts{
				Group:            accessGroup.Group,
				ReviewURLs:       reviewURL,
				RequestReviewer:  reviewerUserObj.Result,
				RequestorEmail:   requestor.Email,
				RequestorSlackID: slackUserID,
				WasReviewed:      true,
			})

			err = n.UpdateMessageBlockForReviewer(ctx, *reqReviewer.Result, slackMsg)

			if err != nil {
				log.Errorw("failed to send slack message", "user", requestor, zap.Error(err))
			}
		}

	}

	// 🚨🚨 TODO
	//
	// Decide on level of noise,
	// Do we want slack webhooks to trigger on every access group review?

	// if HAS_SLACK_WEBHOOKS {
	// Note: propably don't need webhook alerts here...
	// for _, webhook := range n.webhooks {
	// 	_, msg := BuildRequestDetailMessage(RequestDetailMessageOpts{
	// 		Request: accessGroup,
	// 		// RequestArguments: requestArguments,
	// 		HeadingMessage: headingMsg,
	// 	})
	// 	err := webhook.SendWebhookMessage(ctx, msg.Blocks, summary)
	// 	if err != nil {
	// 		log.Errorw("failed to send slack message to webhook channel", "error", err)
	// 	}
	// }
	// }

}
