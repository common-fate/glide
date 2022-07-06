package rulesvc

import "errors"

var (

	// ErrRuleNotFound is returned if a rule with the supplied id already exists.
	ErrRuleIdAlreadyExists = errors.New("access rule id already exists")

	// ErrUserNotAuthorized is returned if the user isn't allowed to complete an action,
	// like reviewing a request.
	ErrUserNotAuthorized = errors.New("user is not authorized to perform this action")

	// ErrProviderNotFound is returned if a matching provider could not be found in the access handler
	ErrProviderNotFound = errors.New("provider not found")

	ErrUnhandledResponseFromAccessHandler = errors.New("access handler returned an unhandled response")

	// ErrAccessRuleAlreadyArchived is returned if an archive request is made for a rule which is already archived
	ErrAccessRuleAlreadyArchived = errors.New("access rule already archived")
)
