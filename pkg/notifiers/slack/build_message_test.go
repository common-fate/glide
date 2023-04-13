package slacknotifier

// func TestBuildRequestMessage(t *testing.T) {
// 	reason := "reason"
// 	start := time.Date(2022, 1, 1, 1, 0, 0, 0, time.UTC)

// 	tests := []struct {
// 		name        string
// 		args        RequestMessageOpts
// 		wantSummary string
// 		wantMsg     string
// 	}{
// 		{
// 			name: "ok",
// 			args: RequestMessageOpts{
// 				Request: access.Request{
// 					ID:          "123",
// 					RequestedBy: "usr_1",
// 					Data: access.RequestData{
// 						Reason: &reason,
// 					},
// 					RequestedTiming: requests.Timing{
// 						Duration:  time.Hour,
// 						StartTime: &start,
// 					},
// 				},
// 				Rule: rule.AccessRule{
// 					Name: "my rule",
// 				},
// 				RequestorEmail: "testuser@example.com",
// 			},
// 			wantSummary: "New request for my rule from testuser@example.com",
// 			wantMsg: `
// {
// 	"replace_original": false,
// 	"delete_original": false,
// 	"metadata": {
// 		"event_type": "",
// 		"event_payload": null
// 	},
// 	"blocks": [
// 		{
// 			"type": "section",
// 			"text": {
// 				"type": "mrkdwn",
// 				"text": "*\u003c|New request for my rule\u003e from testuser@example.com*"
// 			}
// 		},
// 		{
// 			"type": "section",
// 			"fields": [
// 				{
// 					"type": "mrkdwn",
// 					"text": "*When:*\n\u003c!date^1640998800^{date_short_pretty} at {time}|2022-01-01 01:00:00 +0000 UTC\u003e"
// 				},
// 				{
// 					"type": "mrkdwn",
// 					"text": "*Duration:*\n1h0m0s"
// 				},
// 				{
// 					"type": "mrkdwn",
// 					"text": "*Status:*\n"
// 				},
// 				{
// 					"type": "mrkdwn",
// 					"text": "*Request Reason:*\nreason"
// 				}
// 			]
// 		}
// 	]
// }`,
// 		},
// 		{
// 			name:        "doesnt panic if bad data is provided",
// 			args:        RequestMessageOpts{},
// 			wantSummary: "New request for  from ",
// 			wantMsg: `
// {
// 	"replace_original": false,
// 	"delete_original": false,
// 	"metadata": {
// 		"event_type": "",
// 		"event_payload": null
// 	},
// 	"blocks": [
// 		{
// 			"type": "section",
// 			"text": {
// 				"type": "mrkdwn",
// 				"text": "*\u003c|New request for \u003e from *"
// 			}
// 		},
// 		{
// 			"type": "section",
// 			"fields": [
// 				{
// 					"type": "mrkdwn",
// 					"text": "*When:*\nASAP"
// 				},
// 				{
// 					"type": "mrkdwn",
// 					"text": "*Duration:*\n0s"
// 				},
// 				{
// 					"type": "mrkdwn",
// 					"text": "*Status:*\n"
// 				}
// 			]
// 		}
// 	]
// }`,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			gotSummary, gotMsg := BuildRequestReviewMessage(tt.args)
// 			assert.Equal(t, tt.wantSummary, gotSummary)

// 			gotMsgJSON, err := json.MarshalIndent(gotMsg, "", "\t")
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			wantMsg := strings.TrimSpace(tt.wantMsg)
// 			assert.Equal(t, wantMsg, string(gotMsgJSON))
// 		})
// 	}
// }
