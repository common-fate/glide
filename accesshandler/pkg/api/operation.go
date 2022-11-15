package api

import (
	"encoding/json"
	"net/http"

	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/config"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/types"
)

// Call an RPC operation
// (POST /api/v1/grants/{grantId}/operation)
func (a *API) CallOperation(w http.ResponseWriter, r *http.Request, grantId string) {
	ctx := r.Context()
	var b types.OperationRequest
	err := apio.DecodeJSONBody(w, r, &b)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	prov, ok := config.Providers[b.Provider]
	if !ok {
		apio.Error(ctx, w, apio.NewRequestError(&providers.ProviderNotFoundError{Provider: b.Provider}, http.StatusNotFound))
		return
	}

	opprov, ok := prov.Provider.(providers.Operationer)
	if !ok {
		apio.ErrorString(ctx, w, "provider does not implement rpc operations", http.StatusBadRequest)
		return
	}

	ops := opprov.Operations()
	op, ok := ops[b.Operation]
	if !ok {
		apio.ErrorString(ctx, w, "operation not found", http.StatusBadRequest)
		return
	}

	grantArgs, err := json.Marshal(b.With)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	operationArgs, err := json.Marshal(b.With)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	res, err := op.Execute(ctx, providers.OperationOpts{
		Subject:       string(b.Subject),
		GrantArgs:     grantArgs,
		OperationArgs: operationArgs,
		GrantID:       b.RequestId,
	})

	if err != nil {
		apio.ErrorString(ctx, w, err.Error(), http.StatusBadRequest)
		return
	}
	result := types.OperationResponse{
		Data: res,
	}

	apio.JSON(ctx, w, result, http.StatusOK)
}
