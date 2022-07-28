package slacknotifier

import (
	"context"

	"github.com/slack-go/slack"
)

// UpdateMessageBlocks is a utility for updating DMs to users by ID
//
// The message is in Slack message block format.
// DE = support for changing status & removing/disabling block actions
func UpdateMessageBlocks(ctx context.Context, slackClient *slack.Client, userEmail string, message slack.Message) error {
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

	_, _, _, err = slackClient.UpdateMessageContext(ctx, result.Conversation.ID, message.Timestamp, slack.MsgOptionBlocks(message.Blocks.BlockSet...))

	if err != nil {
		return err
	}

	return nil
}
