package slacknotifier

import (
	"testing"
	"time"

	"github.com/common-fate/granted-approvals/pkg/access"
	"github.com/common-fate/granted-approvals/pkg/identity"
	"github.com/common-fate/granted-approvals/pkg/notifiers"
	"github.com/common-fate/granted-approvals/pkg/rule"
	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
)

// test cases for BuildRequestMessage(
func TestBuildRequestMessage(t *testing.T) {

	now := time.Now()

	/*
		Test Cases:
		- General properties from the input match the output (request reason, duration, etc)
		- ActionBlocks are dependent on request status
		- ASAP Requests are handled
	*/
	type want struct {
		summary string
		msg     slack.Msg
	}

	type testcase struct {
		name string
		in   RequestMessageOpts
		want want
	}

	uID1 := "user1"
	uID2 := "user2"
	rID1 := "rule1"

	testcases := []testcase{
		{
			name: "ok",
			in: RequestMessageOpts{
				Request: access.Request{
					ID:          "123",
					Status:      access.PENDING,
					RequestedBy: uID1,
					Rule:        rID1,
					RuleVersion: "1.0.0",
					Data:        access.RequestData{},
				},
				Rule:             rule.AccessRule{},
				RequestorSlackID: "",
				RequestorEmail:   "",
				Reviewer: &identity.User{
					Email: uID2,
				},
				RequestReviewer: &identity.User{
					Email: uID2,
				},
				ReviewURLs: notifiers.ReviewURLs{},
			},
			want: want{
				summary: "New request for  from ",
				msg:     slack.Msg{ClientMsgID: "", Type: "", Channel: "", User: "", Text: "", Timestamp: "", ThreadTimestamp: "", IsStarred: false, PinnedTo: []string(nil), Attachments: []slack.Attachment(nil), Edited: (*slack.Edited)(nil), LastRead: "", Subscribed: false, UnreadCount: 0, SubType: "", Hidden: false, DeletedTimestamp: "", EventTimestamp: "", BotID: "", Username: "", Icons: (*slack.Icon)(nil), BotProfile: (*slack.BotProfile)(nil), Inviter: "", Topic: "", Purpose: "", Name: "", OldName: "", Members: []string(nil), ReplyCount: 0, Replies: []slack.Reply(nil), ParentUserId: "", LatestReply: "", Files: []slack.File(nil), Upload: false, Comment: (*slack.Comment)(nil), ItemType: "", ReplyTo: 0, Team: "", Reactions: []slack.ItemReaction(nil), ResponseType: "", ReplaceOriginal: false, DeleteOriginal: false, Blocks: slack.Blocks{BlockSet: []slack.Block(nil)}},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {

			summ, msg := BuildRequestMessage(tc.in)

			msg.Timestamp = now.Format(time.RFC3339)
			tc.want.msg.Timestamp = now.Format(time.RFC3339)

			assert.Equal(t, msg, tc.want.msg)

			// m1, err := msg.Blocks.MarshalJSON()
			// assert.NoError(t, err)

			// m2, err := tc.want.msg.Blocks.MarshalJSON()
			// assert.NoError(t, err)

			// assert.Equal(t, string(m1), string(m2))
			// m2 := "[{\"type\":\"section\",\"text\":{\"type\":\"mrkdwn\",\"text\":\"*\\u003c|New request for \\u003e from *\"}},{\"type\":\"section\",\"fields\":[{\"type\":\"mrkdwn\",\"text\":\"*When:*\\nASAP\"},{\"type\":\"mrkdwn\",\"text\":\"*Duration:*\\n0s\"},{\"type\":\"mrkdwn\",\"text\":\"*Status:*\\nPending\"}]},{\"type\":\"context\",\"elements\":[{\"type\":\"mrkdwn\",\"text\":\"*Reviewed by* user2 at \\u003c!date^1660032577^{date_short_pretty} at {time}|2022-08-09 16:09:37.011756 +0800 AWST m=+0.007000546\\u003e\"}]},{\"type\":\"actions\",\"block_id\":\"review_actions\",\"elements\":[{\"type\":\"button\",\"text\":{\"type\":\"plain_text\",\"text\":\"Approve\"},\"action_id\":\"approve\",\"value\":\"approve\",\"style\":\"primary\"},{\"type\":\"button\",\"text\":{\"type\":\"plain_text\",\"text\":\"Close Request\"},\"action_id\":\"deny\",\"value\":\"deny\",\"style\":\"danger\"}]}]"

			// assert.Equal(t, tc.want.msg, msg)
			assert.Equal(t, tc.want.summary, summ)
		})
	}
}

// func TestBuildRequestMessage(t *testing.T) {
// 	type args struct {
// 		Request *access.Request
// 	}
// 	tests := []struct{
// 		name string
// 		args args
// 		in RequestMessageOpts
// 		want RequestMessage
// 	}
// }
