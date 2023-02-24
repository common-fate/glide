package requestroutersvc

import "errors"

var ErrCannotRoute error = errors.New("cannot route to deployment for target group")
