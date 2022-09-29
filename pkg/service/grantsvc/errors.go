package grantsvc

import (
	"errors"
	"fmt"

	"github.com/common-fate/granted-approvals/accesshandler/pkg/types"
)

var (
	//ErrGrantInactive is returned when a inactive grant is attempted to be revoked
	ErrGrantInactive = errors.New("only active grants can be revoked")
	// ErrNoGrant is returned when attempting to revoke a request which has no grant yet
	ErrNoGrant = errors.New("request has no grant")
	// ErrNoGrant is returned when attempting to revoke a request which has no grant yet
)

type GrantValidationError struct {
	Validation types.GrantValidation
}

func (e GrantValidationError) Error() string {
	return fmt.Sprintf("validation on grant failed: %v", e.Validation)
}
