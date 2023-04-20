package accesssvc

import (
	"errors"
	"fmt"

	"github.com/common-fate/common-fate/pkg/types"
)

var (
	// ErrNoMatchingGroup is returned during creating a request if the user is not in the correct group,
	// based on the Access Rule the request relates to.
	ErrNoMatchingGroup = errors.New("user was not in a matching group for the access rule")

	// ErrRuleNotFound is returned if we can't find the Access Rule associated with a request.
	ErrRuleNotFound = errors.New("access rule not found")

	// ErrUserNotAuthorized is returned if the user isn't allowed to complete an action,
	// like reviewing a request.
	ErrUserNotAuthorized = errors.New("user is not authorized to perform this action")

	// ErrRequestCannotBeCancelled is returned if the request is not in the pending status
	ErrRequestCannotBeCancelled = errors.New("only pending requests can be cancelled")

	// ErrRequestOverlapsExistingGrant is returned if the request overlaps an existing grant
	ErrRequestOverlapsExistingGrant = errors.New("this request overlaps an existing grant")
)

// InvalidStatusError is returned if a user tries to review a request which wasn't PENDING.
type InvalidStatusError struct {
	Status types.RequestStatus
}

func (e InvalidStatusError) Error() string {
	return fmt.Sprintf("request has invalid status: %s", e.Status)
}
