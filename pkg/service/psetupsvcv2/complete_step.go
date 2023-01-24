package psetupsvcv2

import (
	"context"
	"errors"
	"fmt"

	"github.com/common-fate/common-fate/pkg/providersetupv2"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
)

// InvalidConfigFieldError is returned if the user attempts to set
// config field values which don't exist on the provider.
type InvalidConfigFieldError struct {
	Key string
}

func (e InvalidConfigFieldError) Error() string {
	return fmt.Sprintf("invalid field: %s", e.Key)
}

var ErrInvalidStepIndex = errors.New("invalid step index")

func (s *Service) CompleteStep(ctx context.Context, setupID string, stepIndex int, body types.ProviderSetupStepCompleteRequest) (*providersetupv2.Setup, error) {
	q := storage.GetProviderSetupV2{
		ID: setupID,
	}

	_, err := s.DB.Query(ctx, &q)
	if err == ddb.ErrNoItems {
		return nil, ErrProviderSetupNotFound
	}
	setup := q.Result

	if stepIndex >= len(setup.Steps) {
		return nil, ErrInvalidStepIndex
	}

	setup.Steps[stepIndex].Complete = body.Complete

	if !body.Complete {
		// if the step is marked incomplete, don't update any values.
		// just mark the step as incomplete and then return.
		err = s.DB.Put(ctx, setup)
		if err != nil {
			return nil, err
		}
		return setup, nil
	}

	// @TODO some config stuff maybe

	err = s.DB.Put(ctx, setup)
	if err != nil {
		return nil, err
	}
	return setup, nil

}
