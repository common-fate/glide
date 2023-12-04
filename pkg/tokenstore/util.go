package tokenstore

import (
	"time"

	"golang.org/x/oauth2"
)

func ShouldRefreshToken(token oauth2.Token, now time.Time) bool {
	expiry := token.Expiry

	// if the token expires up to 5 minutes in the future,
	// refresh it now.
	timeToRefresh := now.Add(5 * time.Minute)

	return expiry.Before(timeToRefresh)
}
