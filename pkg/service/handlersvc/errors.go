package handlersvc

import "errors"

var (
	ErrHandlerIdAlreadyExists  = errors.New("handler id already exists")
	ErrInvalidAwsAccountNumber = errors.New("invalid aws account number")
)
