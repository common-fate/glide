package api

import (
	"net/http"

	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/common-fate/pkg/auth"
	"github.com/common-fate/common-fate/pkg/service/accesssvc"
	"github.com/common-fate/common-fate/pkg/types"
)

// Review a request
// (POST /api/v1/access-group/{id}/review)
func (a *API) UserReviewRequest(w http.ResponseWriter, r *http.Request, requestId string, groupId string) {
	ctx := r.Context()
	var reviewRequest types.ReviewRequest
	err := apio.DecodeJSONBody(w, r, &reviewRequest)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	user := auth.UserFromContext(ctx)
	isAdmin := auth.IsAdmin(ctx)

	err = a.Access.Review(ctx, *user, isAdmin, requestId, groupId, reviewRequest)
	if err == accesssvc.ErrAccesGroupNotFoundOrNoAccessToReview {
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusNotFound))
		return
	}
	if err == accesssvc.ErrAccessGroupAlreadyReviewed {
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusBadRequest))
		return
	}
	if err == accesssvc.ErrGroupCannotBeApprovedBecauseItWillOverlapExistingGrants {
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusBadRequest))
		return
	}
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	apio.JSON(ctx, w, nil, http.StatusNoContent)
}
