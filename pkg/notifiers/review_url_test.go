package notifiers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReviewURL(t *testing.T) {
	type testcase struct {
		name    string
		giveURL string
		giveID  string
		want    ReviewURLs
	}

	testcases := []testcase{
		{
			name:    "ok",
			giveURL: "https://grantedtest.com",
			giveID:  "req_123",
			want: ReviewURLs{
				Review:             "https://grantedtest.com/requests/req_123",
				Approve:            "https://grantedtest.com/requests/req_123?action=approve",
				Deny:               "https://grantedtest.com/requests/req_123?action=deny",
				AccessInstructions: "https://grantedtest.com/requests/req_123#access_instructions",
			},
		},
		{
			name:    "with path",
			giveURL: "https://grantedtest.com/prod",
			giveID:  "req_123",
			want: ReviewURLs{
				Review:             "https://grantedtest.com/prod/requests/req_123",
				Approve:            "https://grantedtest.com/prod/requests/req_123?action=approve",
				Deny:               "https://grantedtest.com/prod/requests/req_123?action=deny",
				AccessInstructions: "https://grantedtest.com/prod/requests/req_123#access_instructions",
			},
		},
	}

	for i := range testcases {
		tc := testcases[i]
		t.Run(tc.name, func(t *testing.T) {
			got, err := ReviewURL(tc.giveURL, tc.giveID)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tc.want, got)
		})
	}

}
