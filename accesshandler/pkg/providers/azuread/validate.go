package azuread

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/hashicorp/go-multierror"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

// https://developer.okta.com/docs/reference/error-codes/#E0000007
var oktaErrorCodeNotFound = "E0000007"

// Validate the access against Okta without actually granting it.
func (p *Provider) Validate(ctx context.Context, subject string, args []byte) error {
	var a Args
	err := json.Unmarshal(args, &a)
	if err != nil {
		return err
	}

	// keep a running track of validation errors.
	var result error

	// The user should exist in Okta.
	_, err = p.client.GetUser(ctx, subject)
	if err != nil {
		var oe *okta.Error
		isOktaErr := errors.As(err, &oe)
		if isOktaErr && oe.ErrorCode == oktaErrorCodeNotFound {
			result = multierror.Append(result, &UserNotFoundError{User: subject})
		} else {
			// we got an error we didn't expect so bail out of any further
			// validation, as we may not be authenticated properly to Okta.
			return err
		}
	}

	// The group we are trying to grant access to should exist in Okta.
	_, err = p.client.GetGroup(ctx, a.GroupID)
	if err != nil {
		var oe *okta.Error
		isOktaErr := errors.As(err, &oe)
		if isOktaErr && oe.ErrorCode == oktaErrorCodeNotFound {
			// If we get this error code, the group wasn't found.
			// We use a specific error type for this.
			err = &GroupNotFoundError{Group: a.GroupID}
		}
		// add the error to our list.
		result = multierror.Append(result, err)
	}

	return result
}
