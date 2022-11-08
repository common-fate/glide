package internalidentitysvc

import (
	"errors"
	"fmt"
)

var ErrNotInternal error = errors.New("cannot update group because it is not an internal group")

type UserNotFoundError struct {
	UserID string
}

func (e UserNotFoundError) Error() string {
	return fmt.Sprintf("user %s does not exist", e.UserID)
}
