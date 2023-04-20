package accesssvc

import (
	"context"
)

type CancelRequestOpts struct {
	CancellerID string
	RequestID   string
}

// CancelRequest cancels a request if it is in pending status.
// Returns an error if the request is invalid.
func (s *Service) CancelRequest(ctx context.Context, opts CancelRequestOpts) error {
	// items := []ddb.Keyer{}
	// now := s.Clock.Now()
	// requestGet := storage.GetRequestV2{ID: opts.RequestID}
	// _, err := s.DB.Query(ctx, &requestGet)
	// if err != nil {
	// 	return err
	// }
	// req := requestGet.Result

	// isAllowed := canCancel(opts, *req)
	// if !isAllowed {
	// 	return ErrUserNotAuthorized
	// }

	// q := storage.ListAccessGroups{RequestID: opts.RequestID}
	// _, err = s.DB.Query(ctx, &q)
	// if err != nil {
	// 	return err
	// }
	// accessGroups := q.Result

	// for _, ag := range accessGroups {
	// 	canBeCancelled := isCancellable(ag)
	// 	if !canBeCancelled {
	// 		return ErrRequestCannotBeCancelled
	// 	}

	// 	ag.Status = requests.CANCELLED
	// 	ag.UpdatedAt = now

	// 	items = append(items, &ag)
	// }

	// // Todo: what should happen with grant types here?

	// req.UpdatedAt = now

	// items = append(items, req)

	// return s.DB.PutBatch(ctx, items...)
	return nil
}

// // users can cancel their own requests.
// func canCancel(opts CancelRequestOpts, request requests.Requestv2) bool {
// 	// canceller must be original requestor
// 	return opts.CancellerID == request.RequestedBy.ID
// }

// // A request can be cancelled if
// func isCancellable(AccessGroup access.AccessGroup) bool {
// 	return AccessGroup.Status == types.RequestStatusPENDING || AccessGroup.Status != requests.CANCELLED
// }
