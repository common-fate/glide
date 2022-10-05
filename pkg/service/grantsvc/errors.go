package grantsvc

import (
	"errors"
	"fmt"
)

var (
	//ErrGrantInactive is returned when a inactive grant is attempted to be revoked
	ErrGrantInactive = errors.New("only active grants can be revoked")
	// ErrNoGrant is returned when attempting to revoke a request which has no grant yet
	ErrNoGrant = errors.New("request has no grant")
	// ErrNoGrant is returned when attempting to revoke a request which has no grant yet
)

// GrantValidationError is returned if grantValidation fails
type GrantValidationError struct {
	ValidationFailureMsg string
}

func (e GrantValidationError) Error() string {
	return fmt.Sprintf("validation failed:\n%s", e.ValidationFailureMsg)
}
