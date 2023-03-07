package requestroutersvc

import "errors"

var ErrCannotRoute error = errors.New("cannot route to a handler because all routes for this target group are invalid")
var ErrNoRoutes error = errors.New("no routes exist for this target group")
