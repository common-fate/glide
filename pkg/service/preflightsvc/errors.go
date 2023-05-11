package preflightsvc

import "errors"

var (
	ErrDuplicateTargetIDsRequested         error = errors.New("duplicate target ids were submitted in the request")
	ErrUserNotAuthorisedForRequestedTarget error = errors.New("user in not authorised to access one or more requested targets")
)
