package accesssvc

import (
	"context"

	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/gevent"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/types"
)

type CancelRequestOpts struct {
	CancellerID string
	RequestID   string
}

// CancelRequest cancels a request if it is in pending status.
// Returns an error if the request is invalid.
func (s *Service) CancelRequest(ctx context.Context, opts CancelRequestOpts) error {

	requestGet := storage.GetRequestWithGroupsWithTargets{ID: opts.RequestID}
	_, err := s.DB.Query(ctx, &requestGet)
	if err != nil {
		return err
	}
	req := requestGet.Result.Request

	isAllowed := canCancel(opts, req)
	if !isAllowed {
		return ErrUserNotAuthorized
	}

	isCancellable := isCancellable(*requestGet.Result)
	if !isCancellable {
		return ErrRequestCannotBeCancelled
	}

	return s.EventPutter.Put(ctx, gevent.RequestCancelledInit{
		Request: *requestGet.Result,
	})
}

// // users can cancel their own requests.
func canCancel(opts CancelRequestOpts, request access.Request) bool {
	// canceller must be original requestor
	return opts.CancellerID == request.RequestedBy.ID
}

// // A access group can be cancelled if
func isCancellable(request access.RequestWithGroupsWithTargets) bool {
	var isCancellable bool
	for _, group := range request.Groups {
		isCancellable = group.Status == types.RequestAccessGroupStatusPENDINGAPPROVAL

	}
	return isCancellable
}
