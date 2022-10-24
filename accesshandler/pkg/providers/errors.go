package providers

import "fmt"

type InvalidArgumentError struct {
	Arg string
}

func (e *InvalidArgumentError) Error() string {
	return fmt.Sprintf("argument %s is not valid", e.Arg)
}

type InvalidGroupIDError struct {
	GroupID string
}

func (e *InvalidGroupIDError) Error() string {
	return fmt.Sprintf("groupID %s is not valid for this argument", e.GroupID)
}

type InvalidGroupValueError struct {
	GroupID    string
	GroupValue string
}

func (e *InvalidGroupValueError) Error() string {
	return fmt.Sprintf("groupValue %s is not valid for this groupID %s", e.GroupValue, e.GroupID)
}

type ProviderNotFoundError struct {
	Provider string
}

func (e *ProviderNotFoundError) Error() string {

	return fmt.Sprintf("no provider found matching: %s", e.Provider)
}
