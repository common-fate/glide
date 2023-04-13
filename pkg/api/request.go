package api

import (
	"net/http"

	"github.com/common-fate/apikit/apio"
)

func (a *API) UserRevokeRequest(w http.ResponseWriter, r *http.Request, requestID string) {
	ctx := r.Context()
	// isAdmin := auth.IsAdmin(ctx)
	// u := auth.UserFromContext(ctx)
	// var req access.Request
	// q := storage.GetRequest{ID: requestID}
	// _, err := a.DB.Query(ctx, &q)
	// if err == ddb.ErrNoItems {
	// 	//grant not found return 404
	// 	apio.Error(ctx, w, apio.NewRequestError(errors.New("request not found or you don't have access to it"), http.StatusNotFound))
	// 	return
	// }
	// if err != nil {
	// 	apio.Error(ctx, w, err)
	// 	return
	// }
	// // user can revoke their own request and admins can revoke any request
	// if q.Result.RequestedBy == u.ID || isAdmin {
	// 	req = *q.Result
	// } else { // reviewers can revoke reviewable requests
	// 	q := storage.GetRequestReviewer{RequestID: requestID, ReviewerID: u.Email}
	// 	_, err := a.DB.Query(ctx, &q)
	// 	if err == ddb.ErrNoItems {
	// 		//grant not found return 404
	// 		apio.Error(ctx, w, apio.NewRequestError(errors.New("request not found or you don't have access to it"), http.StatusNotFound))
	// 		return
	// 	}
	// 	if err != nil {
	// 		apio.Error(ctx, w, err)
	// 		return
	// 	}
	// 	req = q.Result.Request
	// }

	// _, err = a.Workflow.Revoke(ctx, req, u.ID, u.Email)
	// if err == workflowsvc.ErrGrantInactive {
	// 	apio.Error(ctx, w, apio.NewRequestError(err, http.StatusBadRequest))
	// 	return
	// }
	// if err == workflowsvc.ErrNoGrant {
	// 	apio.Error(ctx, w, apio.NewRequestError(err, http.StatusBadRequest))
	// 	return
	// }
	// if err != nil {
	// 	apio.Error(ctx, w, err)
	// 	return
	// }

	// analytics.FromContext(ctx).Track(&analytics.RequestRevoked{
	// 	RequestedBy: req.RequestedBy,
	// 	RevokedBy:   u.ID,
	// 	RuleID:      req.Rule,
	// 	Timing:      req.RequestedTiming.ToAnalytics(),
	// 	HasReason:   req.HasReason(),
	// })

	apio.JSON(ctx, w, nil, http.StatusOK)
}
