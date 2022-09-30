package accesssvc

import (
	"context"

	"github.com/common-fate/granted-approvals/pkg/access"
	"github.com/common-fate/granted-approvals/pkg/gevent"
	"github.com/common-fate/granted-approvals/pkg/storage"
	"github.com/common-fate/granted-approvals/pkg/storage/dbupdate"
)

type CancelRequestOpts struct {
	CancellerID string
	RequestID   string
}

// CancelRequest cancels a request if it is in pending status.
// Returns an error if the request is invalid.
func (s *Service) CancelRequest(ctx context.Context, opts CancelRequestOpts) error {
	q := storage.GetRequest{ID: opts.RequestID}
	_, err := s.DB.Query(ctx, &q)
	if err != nil {
		return err
	}
	req := q.Result
	originalStatus := req.Status
	isAllowed := canCancel(opts, *req)
	if !isAllowed {
		return ErrUserNotAuthorized
	}
	canBeCancelled := isCancellable(*req)
	if !canBeCancelled {
		return ErrRequestCannotBeCancelled
	}

	req.Status = access.CANCELLED
	req.UpdatedAt = s.Clock.Now()
	// we need to save the Review, the updated Request in the database.
	items, err := dbupdate.GetUpdateRequestItems(ctx, s.DB, *req)
	if err != nil {
		return err
	}
	// audit log event
	reqEvent := access.NewStatusChangeEvent(req.ID, req.UpdatedAt, &opts.CancellerID, originalStatus, req.Status)

	err = s.EventPutter.Put(ctx, gevent.RequestCancelled{Request: *req})
	// In a future PR we will shift these events out to be triggered by dynamo db streams
	// This will currently put the app in a strange state if this fails
	if err != nil {
		return err
	}

	items = append(items, &reqEvent)
	return s.DB.PutBatch(ctx, items...)
}

// users can cancel their own requests.
func canCancel(opts CancelRequestOpts, request access.Request) bool {
	// canceller must be original requestor
	return opts.CancellerID == request.RequestedBy
}

// A request can be cancelled if
func isCancellable(request access.Request) bool {
	return request.Status == access.PENDING || request.Grant == nil && request.Status != access.CANCELLED
}
