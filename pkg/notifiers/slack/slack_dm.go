package slacknotifier

import (
	"context"

	"github.com/common-fate/common-fate/pkg/gconfig"
	"github.com/pkg/errors"
	"github.com/slack-go/slack"
)

const NotificationsTypeSlack = "slack"

type SlackDirectMessage struct {
	client   *slack.Client
	apiToken gconfig.SecretStringValue
}

func (s *SlackDirectMessage) Config() gconfig.Config {
	return gconfig.Config{
		gconfig.SecretStringField("apiToken", &s.apiToken, "the Slack API token", gconfig.WithNoArgs("/granted/secrets/notifications/slack/token")),
	}
}

func (s *SlackDirectMessage) Init(ctx context.Context) error {
	s.client = slack.New(s.apiToken.Get())
	return nil
}

func (s *SlackDirectMessage) TestConfig(ctx context.Context) error {
	_, err := s.client.GetUsersContext(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to list users while testing slack configuration")
	}
	return nil
}
