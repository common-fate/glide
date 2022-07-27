package slacknotifier

import (
	"context"
	"fmt"
	"time"

	"github.com/slack-go/slack"
)

// UpdateMessageBlocks is a utility for updating DMs to users by ID
//
// The message is in Slack message block format.
// DE = support for changing status & removing/disabling block actions
func UpdateMessageBlocks(ctx context.Context, slackClient *slack.Client, userEmail string, message slack.MsgOption) error {
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

	// 1.  Fetch or find where the existing request-message block relation is (the link)
	//     - Store this in a variable called "blockID" on the RequestReviewer
	//     - Look up this blockId in the database
	//     - add it to an array of blocks to be updated
	// 2.  Hook into this, ensuring you have the message block contents
	// 3.  Update the message block contents with desired values
	// 4.  Feed this into UpdateMessageContext, send the message

	// Tap into logic here for how messages that are being updated
	// Research further
	// https://github.com/jace-ys/bingsoo/blob/25bc364265edc999c2c7f168bc4701b8e107ee5d/pkg/session/vote.go#L63

	test := slack.MsgOptionBlocks()

	// We now want to update the message
	// timestampe for now
	t := time.Now()
	_, ts, _, err := slackClient.UpdateMessageContext(ctx, result.Conversation.ID, t.String(), test)

	if err != nil {
		return err
	}

	// could also handle update logic here....
	// employ ts to run the update
	fmt.Print(ts)

	return nil
}
