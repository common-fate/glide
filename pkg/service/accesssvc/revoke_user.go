package accesssvc

import (
	"context"
	"errors"

	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/storage"
)

type RevokeUserResult struct {
	// RevokedRequestIDs are the IDs of requests that were successfully revoked
	RevokedRequestIDs []string
	// FailedRequestIDs are the IDs of requests that failed to be revoked due to errors
	FailedRequestIDs []string
}

func (s *Service) RevokeUserAccess(ctx context.Context, userID string) (*RevokeUserResult, error) {
	var reqs []access.Request

	// find APPROVED access requests
	q := storage.ListRequestsForUserAndStatus{
		Status: access.APPROVED,
		UserId: userID,
	}

	_, err := s.DB.Query(ctx, &q)
	if err != nil {
		return nil, err
	}

	reqs = append(reqs, q.Result...)

	// find PENDING access requests
	q = storage.ListRequestsForUserAndStatus{
		Status: access.PENDING,
		UserId: userID,
	}

	_, err = s.DB.Query(ctx, &q)
	if err != nil {
		return nil, err
	}

	reqs = append(reqs, q.Result...)

	return nil, errors.New("not implemented")
}
