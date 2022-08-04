package gconfig

import (
	"errors"
	"fmt"
)

var ErrFieldValueMustNotBeNil error = errors.New("field value must not be nil")

type IncorrectArgumentsToSecretPathFuncError struct {
	ExpectedArgs int
	FoundArgs    int
	Key          string
}

func (e IncorrectArgumentsToSecretPathFuncError) Error() string {
	return fmt.Sprintf("secret path function for %s recieved an unexpected number of arguments, expected %d, found %d", e.Key, e.ExpectedArgs, e.FoundArgs)
}
