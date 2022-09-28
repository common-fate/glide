package psetupsvc

import (
	"context"
	"errors"
	"fmt"
	"strings"

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

	// load the current values into this config
	// Skip secrets is set to true because we never need to read the value of a secret in this context.
	// If a value is provided for a secret it will be updated
	err = cfg.Load(ctx, &gconfig.MapLoader{Values: setup.ConfigValues, SkipLoadingSecrets: true})
	if err != nil {
		return nil, err
	}

	// verify that all fields actually correspond to the provider
	for key, value := range body.ConfigValues {
		f, err := cfg.FindFieldByKey(key)
		if err != nil {
			return nil, InvalidConfigFieldError{Key: key}
		}

		// any secret starting with 'awsssm://' is assumed to be an existing reference
		// to a secret and is not set.
		if f.IsSecret() && strings.HasPrefix(value, "awsssm://") {
			continue
		}

		err = f.Set(value)
		if err != nil {
			return nil, err
		}
	}

	newConfig, err := cfg.Dump(ctx, gconfig.SSMDumper{Suffix: s.DeploymentSuffix, SecretPathArgs: []interface{}{setupID}})
	if err != nil {
		return nil, err
	}

	// when using SSMDumper here it returns 'awsssm://' for values which haven't been Set.
	// to work around this, we eliminate empty values from the returned map to avoid overwriting
	// the existing reference to the SSM secret.
	for k := range newConfig {
		if newConfig[k] == "awsssm://" {
			delete(newConfig, k)
		}
	}

	for k, v := range newConfig {
		setup.ConfigValues[k] = v
	}

	err = s.DB.Put(ctx, setup)
	if err != nil {
		return nil, err
	}
	return setup, nil

}
