package psetupsvc

import (
	"context"
	"errors"
	"fmt"

	"github.com/common-fate/ddb"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providerregistry"
	"github.com/common-fate/granted-approvals/pkg/gconfig"
	"github.com/common-fate/granted-approvals/pkg/providersetup"
	"github.com/common-fate/granted-approvals/pkg/storage"
	"github.com/common-fate/granted-approvals/pkg/types"
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

func (s *Service) CompleteStep(ctx context.Context, setupID string, stepIndex int, body types.ProviderSetupStepCompleteRequest) (*providersetup.Setup, error) {
	q := storage.GetProviderSetup{
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

	r := providerregistry.Registry()
	rp, err := r.Lookup(setup.ProviderType, setup.ProviderVersion)
	if err != nil {
		return nil, err
	}
	p := rp.Provider

	var cfg gconfig.Config
	if configer, ok := p.(gconfig.Configer); ok {
		cfg = configer.Config()
	}

	// verify that all fields actually correspond to the provider
	for key, value := range body.ConfigValues {
		_, err = cfg.FindFieldByKey(key)
		if err != nil {
			return nil, InvalidConfigFieldError{Key: key}
		}
		// todo: if the field is secret, it should be written to SSM
		setup.ConfigValues[key] = value
	}

	err = s.DB.Put(ctx, setup)
	if err != nil {
		return nil, err
	}
	return setup, nil

}
