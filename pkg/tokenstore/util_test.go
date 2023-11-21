package tokenstore

import (
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/stretchr/testify/assert"

	"golang.org/x/oauth2"
)

func TestShouldRefresh(t *testing.T) {
	type testcase struct {
		name          string
		token         oauth2.Token
		now           clock.Clock
		shouldRefresh bool
	}

	testcases := []testcase{
		{
			name: "24 hours away passes",
			token: oauth2.Token{
				Expiry: time.Now().Add(time.Hour * 24),
			},
			now:           clock.New(),
			shouldRefresh: false,
		},
		{
			name: "2 minutes away refreshes",
			token: oauth2.Token{
				Expiry: time.Now().Add(time.Minute * 2),
			},
			now: clock.New(),

			shouldRefresh: true,
		},
		{
			name: "5 minutes away refreshes",
			token: oauth2.Token{
				Expiry: time.Now().Add(time.Minute * 5),
			},
			now: clock.New(),

			shouldRefresh: true,
		},
		{
			name: "5 minutes after refreshes",
			token: oauth2.Token{
				Expiry: time.Now().Add(time.Minute * -5),
			},
			now: clock.New(),

			shouldRefresh: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {

			outcome := ShouldRefreshToken(tc.token, tc.now.Now())

			assert.Equal(t, outcome, tc.shouldRefresh)

		})
	}

}
