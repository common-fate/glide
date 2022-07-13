package deploy

import "errors"

var ErrConfigNotExist = errors.New("config does not exist")
var ErrConfigNotNotSetInContext = errors.New("config has not been set in context")
