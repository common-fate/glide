package requestroutersvc

import "errors"

var ErrCannotRoute error = errors.New("cannot route to deployment for target group")
var ErrNoValidRoute error = errors.New("no valid routes found for target group")
