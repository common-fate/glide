package targetsvc

import "errors"

var (

	// ErrRuleNotFound is returned if a rule with the supplied id already exists.
	ErrTargetGroupIdAlreadyExists = errors.New("target group id already exists")

	// ErrProviderNotFound is returned if a matching provider could not be found in the registry
	ErrProviderNotFoundInRegistry = errors.New("provider not found in registry")
)
