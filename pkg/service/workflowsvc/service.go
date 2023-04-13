package workflowsvc

import (
	"context"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/common-fate/common-fate/pkg/gevent"
	"github.com/common-fate/common-fate/pkg/requests"
	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
	"github.com/common-fate/iso8601"
	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
)

//go:generate go run github.com/golang/mock/mockgen -destination=mocks/runtime.go -package=mocks . Runtime
type Runtime interface {
	// isForTargetGroup tells the runtime how to process the request
	// grant is expected to be asyncronous
	Grant(ctx context.Context, grant types.CreateGrant) error
	// isForTargetGroup tells the runtime how to process the request
	// revoke is expected to be syncronous
	Revoke(ctx context.Context, grantID string) error
}

//go:generate go run github.com/golang/mock/mockgen -destination=mocks/eventputter.go -package=mocks . EventPutter
type EventPutter interface {
	Put(ctx context.Context, detail gevent.EventTyper) error
}
type Service struct {
	Runtime  Runtime
	DB       ddb.Storage
	Clk      clock.Clock
	Eventbus EventPutter
}

func (s *Service) Grant(ctx context.Context, access_group requests.AccessGroup, subject string) ([]requests.Grantv2, error) {
	// Contains logic for preparing a grant and emitting events

	now := time.Now()
	items := []ddb.Keyer{}

	if !access_group.AccessRule.Approval.IsRequired() {
		for _, entitlement := range access_group.With {
			createGrant, err := s.prepareCreateGrantRequest(ctx, access_group.TimeConstraints, types.NewRequestID(), subject, access_group.AccessRule, entitlement)
			if err != nil {
				return nil, err
			}
			err = s.Runtime.Grant(ctx, createGrant)
			if err != nil {
				return nil, err
			}

			err = s.Eventbus.Put(ctx, &gevent.GrantCreated{Grant: types.Grant{
				ID:       createGrant.Id,
				Provider: createGrant.Provider,
				End:      createGrant.End.Time,
				Start:    createGrant.Start.Time,
				Status:   types.GrantStatusPENDING,
				Subject:  createGrant.Subject,
				With:     types.Grant_With(createGrant.With),
			}})
			if err != nil {
				return nil, err
			}
			grant := requests.Grantv2{
				ID:          createGrant.Id,
				AccessGroup: access_group.ID,
				Subject:     string(createGrant.Subject),
				Start:       createGrant.Start.Time,
				End:         createGrant.End.Time,
				Status:      types.GrantStatusPENDING,
				With:        types.Grant_With(createGrant.With),
				CreatedAt:   now,
				UpdatedAt:   now,
			}
			access_group.Grants = append(access_group.Grants, grant)
			items = append(items, &access_group)

		}
	}

	err := s.DB.PutBatch(ctx, items...)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Revoke attepmts to syncronously revoke access to a request
// If it is successful, the request is updated in the database, and the updated request is returned from this method
func (s *Service) Revoke(ctx context.Context, request requests.Requestv2, revokerID string, revokerEmail string) (*requests.Requestv2, error) {

	for _, access_group := range request.Groups {
		if access_group.Grants == nil || len(access_group.Grants) == 0 {
			return nil, ErrNoGrant
		}
		for _, grant := range access_group.Grants {

			//Cannot request to revoke/cancel grant if it is not active or pending (state function has been created and executed)
			canRevoke := grant.Status == types.GrantStatusACTIVE || grant.Status == types.GrantStatusPENDING

			if !canRevoke || grant.End.Before(s.Clk.Now()) {
				return nil, ErrGrantInactive
			}

			q := storage.GetAccessRuleCurrent{
				ID: access_group.AccessRule.ID,
			}

			_, err := s.DB.Query(ctx, &q)
			if err != nil {
				return nil, err
			}
			err = s.Runtime.Revoke(ctx, grant.ID)
			if err != nil {
				return nil, err
			}

			// previousStatus := grant.Status
			grant.Status = types.GrantStatusREVOKED
			grant.UpdatedAt = s.Clk.Now()
			// items, err := dbupdate.GetUpdateRequestItems(ctx, s.DB, request)
			// if err != nil {
			// 	return nil, err
			// }

			// //create a request event for audit loggging request change
			// requestEvent := access.NewGrantStatusChangeEvent(request.ID, grant.UpdatedAt, &revokerID, previousStatus, grant.Status)

			// items = append(items, &requestEvent)

			// err = s.DB.PutBatch(ctx, items...)
			// if err != nil {
			// 	return nil, err
			// }
			// // Emit an event for the grant revoke
			// // We have chosen to emit events from the Common Fate app for grant revocation rather than from the access handler because we are using a syncronous API.
			// // All effects from revoking will be implemented in this syncronous api rather than triggered from the events.
			// // So we update the grant status here and save the grant before emitting the event
			// err = s.Eventbus.Put(ctx, gevent.GrantRevoked{Grant: grant.ToAPI(), Actor: revokerID, RevokerEmail: revokerEmail})
			// if err != nil {
			// 	return nil, err
			// }
		}
	}

	return nil, nil
}

// prepareCreateGrantRequest prepares the data for requesting
func (s *Service) prepareCreateGrantRequest(ctx context.Context, requestTiming requests.Timing, requestId string, subject string, accessRule rule.AccessRule, with map[string]string) (types.CreateGrant, error) {

	start, end := requestTiming.GetInterval(requests.WithNow(s.Clk.Now()))

	req := types.CreateGrant{
		Id:       requestId,
		Provider: accessRule.Target.TargetGroupID,
		With: types.CreateGrant_With{
			AdditionalProperties: make(map[string]string),
		},
		Subject: openapi_types.Email(subject),
		Start:   iso8601.New(start),
		End:     iso8601.New(end),
	}

	//todo: rework this to be used safely around the codebase
	for k, v := range with {
		req.With.AdditionalProperties[k] = v
	}
	// for k, v := range request.SelectedWith {
	// 	req.With.AdditionalProperties[k] = v.Value
	// }

	return req, nil
}
