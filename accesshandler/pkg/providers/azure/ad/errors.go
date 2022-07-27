package ad

import (
	"fmt"
)

type UserNotFoundError struct {
	User string
}

func (e *UserNotFoundError) Error() string {
	return fmt.Sprintf("user %s was not found", e.User)
}

type GroupNotFoundError struct {
	Group string
}

func (e *GroupNotFoundError) Error() string {
	return fmt.Sprintf("group %s was not found", e.Group)
}
