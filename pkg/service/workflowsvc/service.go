package workflowsvc

import (
	"context"

	"github.com/benbjohnson/clock"
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/gevent"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
)

// //go:generate go run github.com/golang/mock/mockgen -destination=mocks/runtime.go -package=mocks . Runtime
type Runtime interface {
	// grant is expected to be asyncronous
	Grant(ctx context.Context, access_group access.GroupTarget) error
	// isForTargetGroup tells the runtime how to process the request
	// revoke is expected to be syncronous
	Revoke(ctx context.Context, grantID string) error
}

// //go:generate go run github.com/golang/mock/mockgen -destination=mocks/eventputter.go -package=mocks . EventPutter
type EventPutter interface {
	Put(ctx context.Context, detail gevent.EventTyper) error
}
type Service struct {
	Runtime  Runtime
	DB       ddb.Storage
	Clk      clock.Clock
	Eventbus EventPutter
}

func (s *Service) Grant(ctx context.Context, requestID string, groupID string) ([]access.GroupTarget, error) {
	q := storage.GetRequestGroupWithTargets{
		RequestID: requestID,
		GroupID:   groupID,
	}
	_, err := s.DB.Query(ctx, &q, ddb.ConsistentRead())
	if err != nil {
		return nil, err
	}
	group := q.Result

	start, end := group.Group.GetInterval(access.WithNow(s.Clk.Now()))
	for i, target := range group.Targets {
		target.Grant = &access.Grant{
			Subject: group.Group.RequestedBy.Email,
			Start:   start,
			End:     end,
			Status:  types.RequestAccessGroupTargetStatusAWAITINGSTART,
		}
		err := s.Runtime.Grant(ctx, target)
		if err != nil {
			//override the status here to error
			target.Grant.Status = types.RequestAccessGroupTargetStatusERROR
			evt := gevent.GrantFailed{
				Grant:  target,
				Reason: err.Error(),
			}
			err = s.Eventbus.Put(ctx, evt)
			if err != nil {
				return nil, err
			}
		}

		group.Targets[i] = target

	}
	err = s.DB.PutBatch(ctx, group.DBItems()...)
	if err != nil {
		return nil, err
	}
	return group.Targets, nil
}

// // Revoke attepmts to syncronously revoke access to a request
// // If it is successful, the request is updated in the database, and the updated request is returned from this method
func (s *Service) Revoke(ctx context.Context, requestID string, groupID string, revokerID string, revokerEmail string) error {
	q := storage.GetRequestGroupWithTargets{
		RequestID: requestID,
		GroupID:   groupID,
	}
	_, err := s.DB.Query(ctx, &q, ddb.ConsistentRead())
	if err != nil {
		return err
	}
	group := q.Result
	for _, target := range group.Targets {

		//Cannot request to revoke/cancel grant if it is not active or pending (state function has been created and executed)
		canRevoke := target.Grant.Status == types.RequestAccessGroupTargetStatusACTIVE ||
			target.Grant.Status == types.RequestAccessGroupTargetStatusAWAITINGSTART

		if !canRevoke || target.Grant.End.Before(s.Clk.Now()) {
			return ErrGrantInactive
		}

		err = s.Runtime.Revoke(ctx, target.ID)
		if err != nil {
			return err
		}
		//emit request group revoke event
		err = s.Eventbus.Put(ctx, gevent.GrantRevoked{
			Grant: target,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
