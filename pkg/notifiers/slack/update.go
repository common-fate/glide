package slacknotifier

import (
	"context"

	"github.com/common-fate/granted-approvals/pkg/access"
	"github.com/common-fate/granted-approvals/pkg/storage"
	"github.com/pkg/errors"
	"github.com/slack-go/slack"
)

// UpdateMessageBlocks is a utility for updating DMs to users by ID
//
// The message is in Slack message block format.
// userId is the reviewer
func (n *SlackNotifier) UpdateMessageBlockForReviewer(ctx context.Context, reviewer access.Reviewer, message slack.Message) error {
	if n.directMessageClient != nil {
		q := storage.GetUser{ID: reviewer.ReviewerID}
		_, err := n.DB.Query(ctx, &q)
		if err != nil {
			return errors.Wrap(err, "getting user")
		}
		if reviewer.Notifications.SlackMessageID == nil {
			return errors.New("cannot update message because Notifications.SlackMessageID id is nil")
		}
		message.Timestamp = *reviewer.Notifications.SlackMessageID

		u, err := n.directMessageClient.client.GetUserByEmailContext(ctx, q.Result.Email)
		if err != nil {
			return err
		}
		result, _, _, err := n.directMessageClient.client.OpenConversationContext(ctx, &slack.OpenConversationParameters{
			Users: []string{u.ID},
		})
		if err != nil {
			return err
		}
		_, _, _, err = n.directMessageClient.client.UpdateMessageContext(ctx, result.Conversation.ID, message.Timestamp, slack.MsgOptionBlocks(message.Blocks.BlockSet...))

		if err != nil {
			return err
		}
	}
	return nil
}
