package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/sethvargo/go-retry"
	"github.com/stretchr/testify/assert"
)

type IsActiver interface {
	IsActive(ctx context.Context) (bool, error)
}

// CheckIsProvisioned calls the underlying integration's API to check that access was
// provisioned or not. It returns an error if the access status doesn't match what we
// wanted.
//
// Passing 'want' as true means that we expect access to have been provisioned (i.e.
// the user should be a member of the Okta group).
// Passing 'want' as false means that we expect access to have been revoked (i.e.
// the user should NOT be a member of the Okta group.
//
// For some integrations the API is eventually consistent, and access won't be
// reflected immediately when calling this function. To handle this, CheckIsProvisioned
// uses go-retry to call the API again with a backoff with a maximum duration of 10 seconds.
func CheckIsProvisioned(ctx context.Context, access IsActiver, want bool) error {
	b := retry.NewFibonacci(time.Second)
	b = retry.WithMaxDuration(time.Second*10, b)

	return retry.Do(ctx, b, func(ctx context.Context) error {
		exists, err := access.IsActive(ctx)
		if err != nil {
			return err
		}

		if exists != want {
			return retry.RetryableError(fmt.Errorf("IsProvisioned: wanted %v but got %v", want, exists))
		}
		return nil
	})
}

// assertAccessError checks the access error against the expected validation error for the test case.
// If we expect validation to fail, granting and revoking access should also fail and accessErr should be non-nil.
// If we expect validation to succeed, granting and revoking access should also succeed and accessErr should be nil.
func AssertAccessError(t *testing.T, wantErr error, accessErr error, msg string) {
	if wantErr != nil {
		// we have an expected validation error, so we expect access to fail.
		assert.NotNilf(t, accessErr, "%s error: expected access to fail, but the returned error was nil", msg)
	} else {
		// we don't have an expected validation error, so we expect access to succeed.
		assert.Nilf(t, accessErr, "%s error: expected access to succeed, but got an error", msg)
	}
}
