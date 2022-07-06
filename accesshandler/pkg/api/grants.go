package api

import (
	"net/http"

	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/types"
)

// List Grants
// (GET /api/v1/grants)
func (a *API) GetGrants(w http.ResponseWriter, r *http.Request) {

}

// Create Grant
// (POST /api/v1/grants)
func (a *API) PostGrants(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var b types.CreateGrant

	grant, err := func() (*types.Grant, error) {
		err := apio.DecodeJSONBody(w, r, &b)
		if err != nil {
			return nil, err
		}

		g, err := b.Validate(ctx, a.Clock.Now())
		if err != nil {
			// return the error details to the client if validation failed
			return nil, &apio.APIError{
				Err:    err,
				Status: http.StatusBadRequest,
			}
		}

		return a.runtime.CreateGrant(ctx, *g)
	}()

	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	res := types.GrantResponse{
		Grant: grant,
	}

	apio.JSON(ctx, w, res, http.StatusCreated)
}

// Revoke grant
// (POST /api/v1/grants/{grantId}/revoke)
func (a *API) PostGrantsRevoke(w http.ResponseWriter, r *http.Request, grantId string) {
	ctx := r.Context()

	g, err := a.runtime.RevokeGrant(ctx, grantId)

	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	res := types.GrantResponse{
		Grant: g,
	}

	apio.JSON(ctx, w, res, http.StatusOK)
}
