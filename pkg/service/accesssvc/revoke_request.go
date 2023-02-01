package accesssvc

import (
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/identity"
)

type RevokeRequestOpts struct {
	Request        access.Request
	Revoker        identity.User
	RevokerIsAdmin bool
}

type RevokeRequestResult struct{}

// func (s *Service) RevokeRequest(ctx context.Context, in RevokeRequestOpts) (*RevokeRequestResult, error) {
// 	var req access.Request
// 	// user can revoke their own request and admins can revoke any request
// 	if in.Request.RequestedBy == in.Revoker.ID || in.RevokerIsAdmin {
// 		req = in.Request
// 	} else {
// 		// reviewers can revoke reviewable requests
// 		q := storage.GetRequestReviewer{RequestID: in.Request.ID, ReviewerID: in.Revoker.ID}
// 		_, err := s.DB.Query(ctx, &q)
// 		if err == ddb.ErrNoItems {
// 			// reviewer not found
// 			return nil, &apio.APIError{Err: errors.New("request not found or you don't have access to it"), Status: http.StatusNotFound}
// 		}
// 		if err != nil {
// 			return nil, err
// 		}
// 		req = q.Result.Request
// 	}

// 	_, err = a.Granter.RevokeGrant(ctx, grantsvc.RevokeGrantOpts{Request: req, RevokerID: uid})
// 	if err == grantsvc.ErrGrantInactive {
// 		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusBadRequest))
// 		return
// 	}
// 	if err == grantsvc.ErrNoGrant {
// 		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusBadRequest))
// 		return
// 	}
// 	if err != nil {
// 		apio.Error(ctx, w, err)
// 		return
// 	}
// }
