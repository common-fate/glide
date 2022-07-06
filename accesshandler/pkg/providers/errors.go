package providers

import "fmt"

type InvalidArgumentError struct {
	Arg string
}

func (e *InvalidArgumentError) Error() string {
	return fmt.Sprintf("argument %s is not valid", e.Arg)
}

type ProviderNotFoundError struct {
	Provider string
}

func (e *ProviderNotFoundError) Error() string {

	return fmt.Sprintf("no provider found matching: %s", e.Provider)
}
