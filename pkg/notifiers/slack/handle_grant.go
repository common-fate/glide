package slacknotifier

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/common-fate/common-fate/pkg/gevent"
	"github.com/common-fate/common-fate/pkg/notifiers"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
)

func (n *SlackNotifier) HandleGrantEvent(ctx context.Context, log *zap.SugaredLogger, event events.CloudWatchEvent) error {

	var grantEvent gevent.GrantEventPayload
	err := json.Unmarshal(event.Detail, &grantEvent)
	if err != nil {
		return err
	}

	gq := storage.GetRequestV2{ID: grantEvent.Grant.ID}
	_, err = n.DB.Query(ctx, &gq)
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
		msg = fmt.Sprintf("Your access to *%s* is now active.", gq.Result.ID)
		accessory = &slack.Accessory{
			ButtonElement: &slack.ButtonBlockElement{
				Type: slack.METButton,
				Text: slack.NewTextBlockObject(slack.PlainTextType, "Access Instructions", true, false),
				URL:  reviewURL.AccessInstructions,
			},
		}
		fallback = fmt.Sprintf("Your access to %s is now active.", gq.Result.ID)
	case gevent.GrantFailedType:
		msg = fmt.Sprintf("We've had an issue trying to provision or clean up your access to *%s*. We'll keep trying, but if you urgently need access to the role please contact your cloud administrator.", gq.Result.ID)
		fallback = fmt.Sprintf("We've had an issue with your access to %s", gq.Result.ID)
	case gevent.GrantRevokedType:
		msg = fmt.Sprintf("Your access to *%s* has been cancelled by your administrator. Please contact your cloud administrator for more information.", gq.Result.ID)
		fallback = fmt.Sprintf("Your access to %s has been cancelled by your administrator", gq.Result.ID)
	default:
		zap.S().Infow("unhandled grant event", "detailType", event.DetailType)
	}
	if msg != "" {

		_, err = SendMessage(ctx, n.directMessageClient.client, gq.Result.RequestedBy.Email, msg, fallback, accessory)
		return err
	}
	return nil
}
