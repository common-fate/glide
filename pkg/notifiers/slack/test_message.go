package slacknotifier

import (
	"context"
	"encoding/json"

	"github.com/common-fate/granted-approvals/pkg/config"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/slack-go/slack"
)

// SendTestMessage is a helper used for customers to test their slack integration settings
// it expects the configuration json as an input which is parsed and loaded from ssm
func SendTestMessage(ctx context.Context, email string, slackConfig []byte) error {
	var s deploy.SlackConfig
	err := json.Unmarshal(slackConfig, &s)
	if err != nil {
		return err
	}
	err = config.LoadAndReplaceSSMValues(ctx, &s)
	if err != nil {
		panic(err)
	}

	slackClient := slack.New(s.APIToken)
	_, err = SendMessage(ctx, slackClient, email, "slack integration test", "slack integration test")
	if err != nil {
		return err
	}
	return nil
}
