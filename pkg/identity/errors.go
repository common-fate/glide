package identity

import "fmt"

type UserNotFoundError struct {
	User string
}

func (e UserNotFoundError) Error() string {
	return fmt.Sprintf("user %s not found", e.User)
}
