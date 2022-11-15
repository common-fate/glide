package api

import (
	"encoding/json"
	"net/http"

	"github.com/common-fate/apikit/apio"
	ahtypes "github.com/common-fate/granted-approvals/accesshandler/pkg/types"
	"github.com/common-fate/granted-approvals/pkg/auth"
	"github.com/common-fate/granted-approvals/pkg/storage"
	"github.com/common-fate/granted-approvals/pkg/types"
	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/pkg/errors"
)

// Call an RPC operation associated with an Access Request
// (POST /api/v1/requests/{requestId}/operation)
func (a *API) CallRequestOperation(w http.ResponseWriter, r *http.Request, requestId string) {
	ctx := r.Context()

	var b types.OperationRequest
	err := apio.DecodeJSONBody(w, r, &b)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	// get user from context
	uid := auth.UserIDFromContext(ctx)
	q := storage.GetRequest{ID: requestId}
	_, err = a.DB.Query(ctx, &q)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	if q.Result.RequestedBy != uid {
		// not authorised
		apio.ErrorString(ctx, w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	opres, err := a.AccessHandlerClient.CallOperationWithResponse(ctx, q.Result.ID, ahtypes.CallOperationJSONRequestBody{
		Operation:     b.Operation,
		OperationArgs: b.Args,
		Provider:      q.Result.Grant.Provider,
		RequestId:     q.Result.ID,
		Subject:       openapi_types.Email(q.Result.Grant.Subject),
		With:          q.Result.Grant.With.AdditionalProperties,
	})
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	if opres.StatusCode() != http.StatusOK {
		var errres types.ErrorResponse
		err = json.Unmarshal(opres.Body, &errres)
		if err != nil {
			apio.Error(ctx, w, errors.Wrap(err, "unmarshalling error"))
			return
		}
		apio.ErrorString(ctx, w, errres.Error, http.StatusBadRequest)
	}

	res := types.OperationResponse{
		Data: opres.JSON200.Data,
	}

	apio.JSON(ctx, w, res, http.StatusOK)
}
