package sso

import "fmt"

type PermissionSetNotFoundErr struct {
	PermissionSet string
	// the underlying AWS error
	AWSErr error
}

func (e *PermissionSetNotFoundErr) Error() string {
	return fmt.Sprintf("permission set %s was not found or you don't have access to it", e.PermissionSet)
}

type UserNotFoundError struct {
	Email string
}

func (e *UserNotFoundError) Error() string {
	return fmt.Sprintf("could not find user %s in AWS SSO", e.Email)
}

type AccountNotFoundError struct {
	AccountID string
}

func (e *AccountNotFoundError) Error() string {
	return fmt.Sprintf("AWS account %s does not exist in your organization", e.AccountID)
}
