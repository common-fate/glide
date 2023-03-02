package targetsvc

import "errors"

var (

	// ErrRuleNotFound is returned if a rule with the supplied id already exists.
	ErrTargetGroupIdAlreadyExists = errors.New("target group id already exists")

	// ErrProviderNotFound is returned if a matching provider could not be found in the registry
	ErrProviderNotFoundInRegistry = errors.New("provider not found in registry")

	// ErrKindIsRequired is returned if a if the kind was not provided
	ErrKindIsRequired = errors.New("kind is required, provider id should be in the format `common-fate/aws@v0.1.0/Kind`")

	// ErrProviderNotFound is returned if a matching provider could not be found in the registry
	ErrProviderDoesNotImplementKind = errors.New("provider does not implement the kind")
)
