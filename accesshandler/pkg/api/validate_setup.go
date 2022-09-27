package api

import (
	"net/http"
	"sync"

	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/config"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/diagnostics"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providerregistry"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/types"
	"github.com/common-fate/granted-approvals/pkg/gconfig"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

// Validate an Access Provider's settings
// (POST /api/v1/setup/validate)
func (a *API) ValidateSetup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var b types.ValidateRequest
	err := apio.DecodeJSONBody(w, r, &b)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	// look up the provider in the registry
	rp, err := providerregistry.Registry().LookupByUses(b.Uses)
	if err != nil {
		logger.Get(ctx).Error("error looking up provider", zap.Error(err))
		apio.ErrorString(ctx, w, "error looking up provider", http.StatusBadRequest)
		return
	}

	p := rp.Provider
	cv, ok := p.(providers.ConfigValidator)
	if !ok {
		// show a success message, but note that validation has been skipped because the provider doesn't support it.
		res := types.ValidateResponse{
			Validations: []types.ProviderConfigValidation{
				{
					Id:     "no-validation",
					Name:   "Validation has been skipped",
					Status: types.ProviderConfigValidationStatusSUCCESS,
					Logs: []types.Log{
						{Level: types.LogLevelWARNING, Msg: "This Access Provider doesn't support config validation."},
					},
				},
			},
		}

		apio.JSON(ctx, w, res, http.StatusOK)
		return
	}
	validations := cv.ValidateConfig()
	res := types.ValidateResponse{}
	var mu sync.Mutex
	handleResults := func(key string, value providers.ConfigValidationStep, logs diagnostics.Logs) {
		mu.Lock()
		defer mu.Unlock()

		result := types.ProviderConfigValidation{
			Id:              key,
			Name:            value.Name,
			FieldsValidated: value.FieldsValidated,
		}

		if logs.HasSucceeded() {
			result.Status = types.ProviderConfigValidationStatusSUCCESS
		} else {
			result.Status = types.ProviderConfigValidationStatusERROR
		}

		for _, l := range logs {
			result.Logs = append(result.Logs, types.Log{
				Level: types.LogLevel(l.Level),
				Msg:   l.Msg,
			})
		}

		res.Validations = append(res.Validations, result)
	}

	// We will likely encounter an error while initialising the provider internals if some of the values are completely wrong.
	// for example a role ARN being invalid will fail when testing assumed role credentials.
	err = config.SetupProvider(ctx, p, &gconfig.MapLoader{Values: b.With})
	if err != nil {
		logger.Get(ctx).Error("error setting up provider", zap.Error(err))
		// We set the initialisation error on all tests and return to the client
		// because it's difficult for us know know which field exactly caused the error as our init methods on providers are mostly out of our control
		for key, val := range validations {
			handleResults(key, val, diagnostics.Error(err))
		}
		apio.JSON(ctx, w, res, http.StatusOK)
		return
	}

	// if there were no initialisation errors, we can run the validations!
	g, gctx := errgroup.WithContext(ctx)
	for key, val := range validations {
		k := key
		v := val
		g.Go(func() error {
			logs := v.Run(gctx)
			handleResults(k, v, logs)
			return nil
		})
	}

	err = g.Wait()
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	apio.JSON(ctx, w, res, http.StatusOK)
}
