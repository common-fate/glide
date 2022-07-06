package slacknotifier

import (
	"context"

	"github.com/common-fate/granted-approvals/pkg/storage"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
)

// SendMessageBlocks is a utility for sending DMs to users by ID
//
// The message is in Slack message block format.
// The summary must be plaintext and is used as the fallback
// message in Slack notifications.
func SendMessageBlocks(ctx context.Context, slackClient *slack.Client, userEmail string, message slack.Message, summary string) error {
	u, err := slackClient.GetUserByEmailContext(ctx, userEmail)
	if err != nil {
		return err
	}
	result, _, _, err := slackClient.OpenConversationContext(ctx, &slack.OpenConversationParameters{
		Users: []string{u.ID},
	})
	if err != nil {
		return err
	}
	_, _, _, err = slackClient.SendMessageContext(ctx, result.Conversation.ID, slack.MsgOptionBlocks(message.Blocks.BlockSet...))
	return err
}

// SendMessage is a utility for sending DMs to users by ID
//
// The message may be markdown formatted. The summary must be plaintext and is used as the fallback
// message in Slack notifications.
func SendMessage(ctx context.Context, slackClient *slack.Client, userID, message, summary string) error {
	u, err := slackClient.GetUserByEmailContext(ctx, userID)
	if err != nil {
		return err
	}
	result, _, _, err := slackClient.OpenConversationContext(ctx, &slack.OpenConversationParameters{
		Users: []string{u.ID},
	})
	if err != nil {
		return err
	}
	block := slack.NewTextBlockObject("mrkdwn", message, false, false)
	msgBlock := slack.NewSectionBlock(block, nil, nil)
	_, _, _, err = slackClient.SendMessageContext(ctx, result.Conversation.ID, slack.MsgOptionBlocks(msgBlock), slack.MsgOptionText(summary, false))
	return err
}

// SendDMWithLogOnError attempts to fetch a user from cognito to get their email, then tries to send them a message in slack
//
// This will log any errors and continue
func (n *Notifier) SendDMWithLogOnError(ctx context.Context, slackClient *slack.Client, log *zap.SugaredLogger, userId, msg, fallback string) {
	userQuery := storage.GetUser{ID: userId}
	_, err := n.DB.Query(ctx, &userQuery)
	if err != nil {
		log.Errorw("Failed to fetch user by id while trying to send message in slack", "uid", userId, "error", err)
		return
	}
	if err := SendMessage(ctx, slackClient, userQuery.Result.Email, msg, fallback); err != nil {
		log.Errorw("Failed to send direct message", "email", userQuery.Result.Email, "msg", msg, "error", err)
	}
}
