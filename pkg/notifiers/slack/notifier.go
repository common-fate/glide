package slacknotifier

import (
	"context"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/common-fate/ddb"
	"github.com/common-fate/granted-approvals/pkg/gconfig"
	"github.com/pkg/errors"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
)

const NotificationsTypeSlack = "slack"

// Notifier provides handler methods for sending notifications to slack based on events
type SlackNotifier struct {
	DB          ddb.Storage
	FrontendURL string
	client      *slack.Client
	apiToken    gconfig.SecretStringValue
}

func (s *SlackNotifier) Config() gconfig.Config {
	return gconfig.Config{
		Fields: []*gconfig.Field{
			gconfig.SecretStringField("apiToken", &s.apiToken, "the Slack API token", gconfig.WithNoArgs("/granted/secrets/notifications/slack/token")),
		},
	}
}

func (s *SlackNotifier) Init(ctx context.Context) error {
	s.client = slack.New(s.apiToken.Get())
	return nil
}
func (s *SlackNotifier) TestConfig(ctx context.Context) error {
	_, err := s.client.GetUsersContext(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to list users while testing slack configuration")
	}
	return nil
}
func (n *SlackNotifier) HandleEvent(ctx context.Context, event events.CloudWatchEvent) (err error) {
	log := zap.S()

	log.Infow("received event", "event", event)

	if strings.HasPrefix(event.DetailType, "grant") {
		err = n.HandleGrantEvent(ctx, log, event)
		if err != nil {
			return err
		}
	} else if strings.HasPrefix(event.DetailType, "request") {
		err = n.HandleRequestEvent(ctx, log, event)
		if err != nil {
			return err
		}
	} else {
		log.Info("ignoring unhandled event type")
	}
	return nil
}
