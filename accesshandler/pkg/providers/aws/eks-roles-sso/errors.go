package eksrolessso

import "fmt"

type UserNotFoundError struct {
	Email string
}

func (e *UserNotFoundError) Error() string {
	return fmt.Sprintf("could not find user %s in AWS SSO", e.Email)
}
