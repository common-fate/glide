package accesssvc

import (
	"context"

	"github.com/common-fate/common-fate/pkg/requests"
	"github.com/common-fate/common-fate/pkg/storage"
)

type CancelRequestOpts struct {
	CancellerID string
	RequestID   string
}

// CancelRequest cancels a request if it is in pending status.
// Returns an error if the request is invalid.
func (s *Service) CancelRequest(ctx context.Context, opts CancelRequestOpts) error {
	q := storage.GetRequestV2{ID: opts.RequestID}
	_, err := s.DB.Query(ctx, &q)
	if err != nil {
		return err
	}
	req := q.Result

	isAllowed := canCancel(opts, *req)
	if !isAllowed {
		return ErrUserNotAuthorized
	}
	// canBeCancelled := isCancellable(*req)
	// if !canBeCancelled {
	// 	return ErrRequestCannotBeCancelled
	// }

	for _, access_group := range req.Groups {
		access_group.Status = requests.CANCELLED
	}
	req.UpdatedAt = s.Clock.Now()
	// we need to save the Review, the updated Request in the database.
	// items, err := dbupdate.GetUpdateRequestItems(ctx, s.DB, *req)
	// if err != nil {
	// 	return err
	// }

	// return s.DB.PutBatch(ctx, items...)
	return nil
}

// users can cancel their own requests.
func canCancel(opts CancelRequestOpts, request requests.Requestv2) bool {
	// canceller must be original requestor
	return opts.CancellerID == request.RequestedBy.ID
}

// A request can be cancelled if
// func isCancellable(request access.Request) bool {
// 	return request.Status == access.PENDING || request.Grant == nil && request.Status != access.CANCELLED
// }
