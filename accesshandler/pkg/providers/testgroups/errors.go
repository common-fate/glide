package testgroups

import "fmt"

type GroupNotFoundError struct {
	Group string
}

func (e *GroupNotFoundError) Error() string {
	return fmt.Sprintf("group %s was not found", e.Group)
}
