package providers

import "fmt"

type InvalidArgumentError struct {
	Arg string
}

func (e *InvalidArgumentError) Error() string {
	return fmt.Sprintf("argument %s is not valid", e.Arg)
}

type InvalidFilterIdError struct {
	FilterId string
}

func (e *InvalidFilterIdError) Error() string {
	return fmt.Sprintf("argument doesn't support %s filterId", e.FilterId)
}

type ProviderNotFoundError struct {
	Provider string
}

func (e *ProviderNotFoundError) Error() string {

	return fmt.Sprintf("no provider found matching: %s", e.Provider)
}
