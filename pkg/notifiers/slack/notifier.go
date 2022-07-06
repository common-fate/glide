package slacknotifier

import (
	"context"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/common-fate/ddb"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/common-fate/granted-approvals/pkg/types"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
)

type IdentityProvider interface {
	GetUserBySub(ctx context.Context, sub string) (*types.User, error)
}

// Notifier provides handler methods for sending notifications to slack based on events
type Notifier struct {
	DB          ddb.Storage
	FrontendURL string
	SlackConfig deploy.Slack
}

func (n *Notifier) HandleEvent(ctx context.Context, event events.CloudWatchEvent) (err error) {
	log := zap.S()

	log.Infow("received event", "event", event)

	slackClient := slack.New(n.SlackConfig.APIToken)

	if strings.HasPrefix(event.DetailType, "grant") {
		err = n.HandleGrantEvent(ctx, log, slackClient, event)
		if err != nil {
			return err
		}
	} else if strings.HasPrefix(event.DetailType, "request") {
		err = n.HandleRequestEvent(ctx, log, slackClient, event)
		if err != nil {
			return err
		}
	} else {
		log.Info("ignoring unhandled event type")
	}
	return nil
}
