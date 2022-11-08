package slacknotifier

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/common-fate/granted-approvals/pkg/gevent"
	"github.com/common-fate/granted-approvals/pkg/notifiers"
	"github.com/common-fate/granted-approvals/pkg/storage"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
)

func (n *SlackNotifier) HandleGrantEvent(ctx context.Context, log *zap.SugaredLogger, event events.CloudWatchEvent) error {

	var grantEvent gevent.GrantEventPayload
	err := json.Unmarshal(event.Detail, &grantEvent)
	if err != nil {
		return err
	}

	gq := storage.GetRequest{ID: grantEvent.Grant.ID}
	_, err = n.DB.Query(ctx, &gq)
	if err != nil {
		return err
	}
	rq := storage.GetAccessRuleVersion{ID: gq.Result.Rule, VersionID: gq.Result.RuleVersion}
	_, err = n.DB.Query(ctx, &rq)
	if err != nil {
		return err
	}
	var msg string
	var fallback string
	var accessory *slack.Accessory

	reviewURL, err := notifiers.ReviewURL(n.FrontendURL, gq.Result.ID)

	if err != nil {
		return err
	}

	// get the message text based on the event type
	switch event.DetailType {
	case gevent.GrantActivatedType:
		msg = fmt.Sprintf("Your access to *%s* is now active.", rq.Result.Name)
		accessory = &slack.Accessory{
			ButtonElement: &slack.ButtonBlockElement{
				Type: slack.PlainTextType,
				Text: slack.NewTextBlockObject(slack.PlainTextType, "Access Instructions", true, false),
				URL:  reviewURL.AccessInstructions,
			},
		}
		fallback = fmt.Sprintf("Your access to %s is now active.", rq.Result.Name)
	case gevent.GrantExpiredType:
		msg = fmt.Sprintf("Your access to *%s* has now expired. If you still need access you can send another request using Granted Approvals.", rq.Result.Name)
		fallback = fmt.Sprintf("Your access to %s has now expired.", rq.Result.Name)
	case gevent.GrantFailedType:
		msg = fmt.Sprintf("We've had an issue trying to provision or clean up your access to *%s*. We'll keep trying, but if you urgently need access to the role please contact your cloud administrator.", rq.Result.Name)
		fallback = fmt.Sprintf("We've had an issue with your access to %s", rq.Result.Name)
	case gevent.GrantRevokedType:
		msg = fmt.Sprintf("Your access to *%s* has been cancelled by your administrator. Please contact your cloud administrator for more information.", rq.Result.Name)
		fallback = fmt.Sprintf("Your access to %s has been cancelled by your administrator", rq.Result.Name)
	default:
		zap.S().Infow("unhandled grant event", "detailType", event.DetailType)
	}
	if msg != "" {
		_, err = SendMessage(ctx, n.directMessageClient.client, gq.Result.Grant.Subject, msg, fallback, accessory)
		return err
	}
	return nil
}
