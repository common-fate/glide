package types

import (
	"context"
	"errors"
	"time"
)

// ValidCreateGrant is a grant which has been validated.
//
// calling .Validate() on a CreateGrant will give you a
// ValidCreateGrant object if validation passes.
//
// We use ValidCreateGrant wherever possible in the codebase
// to ensure that only validated grants are provisioned.
type ValidCreateGrant struct {
	CreateGrant
}

type ErrInvalidGrantTime struct {
	Msg string
}

func (e ErrInvalidGrantTime) Error() string {
	return e.Msg
}

var ErrInvalidGrantID = errors.New("invalid grant id")

// Validate a grant. If a grant is valid, a ValidGrant object is returned.
//
// Currently this doesn't do any validation of the `with` field against our providers.
func (cg *CreateGrant) Validate(ctx context.Context, now time.Time) (*ValidCreateGrant, error) {
	if cg.Id == "" {
		return nil, ErrInvalidGrantID
	}
	if cg.Start.Equal(cg.End.Time) {
		return nil, ErrInvalidGrantTime{"grant start and end time cannot be equal"}
	}

	if cg.Start.After(cg.End.Time) {
		return nil, ErrInvalidGrantTime{"grant start time must be earlier than end time"}
	}

	if cg.End.Before(now) {
		return nil, ErrInvalidGrantTime{"grant finish time is in the past"}
	}

	return &ValidCreateGrant{*cg}, nil
}
