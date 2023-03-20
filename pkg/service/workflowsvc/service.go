package workflowsvc

import (
	"context"

	"github.com/benbjohnson/clock"
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/gevent"
	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/storage/dbupdate"
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

func (s *Service) Grant(ctx context.Context, request access.Request, accessRule rule.AccessRule) (*access.Grant, error) {
	// Contains logic for preparing a grant and emitting events
	createGrant, err := s.prepareCreateGrantRequest(ctx, request, accessRule)
	if err != nil {
		return nil, err
	}
	err = s.Runtime.Grant(ctx, createGrant)
	if err != nil {
		return nil, err
	}

	// @TODO because v1 access providers sends this event in the rest API method
	// We have to skip sending here here unless it is for a target group
	// in which case the event can be sent from here
	// There is still some issues with a race condition here, for asap requests, the grant could start in step functions before this event it sent.
	// or they could arrive at the same times, if instead this event is produced by the step functions or local mock workflow, then the race condition won't exist
	// SHIFT THIS TO THE STEP FUNCTION
	if accessRule.Target.IsForTargetGroup() {
		err = s.Eventbus.Put(ctx, &gevent.GrantCreated{Grant: types.Grant{
			ID:       createGrant.Id,
			Provider: createGrant.Provider,
			End:      createGrant.End,
			Start:    createGrant.Start,
			Status:   types.GrantStatusPENDING,
			Subject:  createGrant.Subject,
			With:     types.Grant_With(createGrant.With),
		}})
		if err != nil {
			return nil, err
		}
	}
	now := s.Clk.Now()
	return &access.Grant{
		Provider:  createGrant.Provider,
		Subject:   string(createGrant.Subject),
		Start:     createGrant.Start,
		End:       createGrant.End,
		Status:    types.GrantStatusPENDING,
		With:      types.Grant_With(createGrant.With),
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// Revoke attepmts to syncronously revoke access to a request
// If it is successful, the request is updated in the database, and the updated request is returned from this method
func (s *Service) Revoke(ctx context.Context, request access.Request, revokerID string, revokerEmail string) (*access.Request, error) {
	if request.Grant == nil {
		return nil, ErrNoGrant
	}
	//Cannot request to revoke/cancel grant if it is not active or pending (state function has been created and executed)
	canRevoke := request.Grant.Status == types.GrantStatusACTIVE || request.Grant.Status == types.GrantStatusPENDING

	if !canRevoke || request.Grant.End.Before(s.Clk.Now()) {
		return nil, ErrGrantInactive
	}

	q := storage.GetAccessRuleVersion{
		ID:        request.Rule,
		VersionID: request.RuleVersion,
	}

	_, err := s.DB.Query(ctx, &q)
	if err != nil {
		return nil, err
	}
	err = s.Runtime.Revoke(ctx, request.ID)
	if err != nil {
		return nil, err
	}

	previousStatus := request.Grant.Status
	request.Grant.Status = types.GrantStatusREVOKED
	request.Grant.UpdatedAt = s.Clk.Now()
	items, err := dbupdate.GetUpdateRequestItems(ctx, s.DB, request)
	if err != nil {
		return nil, err
	}

	//create a request event for audit loggging request change
	requestEvent := access.NewGrantStatusChangeEvent(request.ID, request.Grant.UpdatedAt, &revokerID, previousStatus, request.Grant.Status)

	items = append(items, &requestEvent)

	err = s.DB.PutBatch(ctx, items...)
	if err != nil {
		return nil, err
	}
	// Emit an event for the grant revoke
	// We have chosen to emit events from the Common Fate app for grant revocation rather than from the access handler because we are using a syncronous API.
	// All effects from revoking will be implemented in this syncronous api rather than triggered from the events.
	// So we update the grant status here and save the grant before emitting the event
	err = s.Eventbus.Put(ctx, gevent.GrantRevoked{Grant: request.Grant.ToAPI(), Actor: revokerID, RevokerEmail: revokerEmail})
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// prepareCreateGrantRequest prepares the data for requesting
func (s *Service) prepareCreateGrantRequest(ctx context.Context, request access.Request, accessRule rule.AccessRule) (types.CreateGrant, error) {
	q := &storage.GetUser{
		ID: request.RequestedBy,
	}
	_, err := s.DB.Query(ctx, q)
	if err != nil {
		return types.CreateGrant{}, err
	}

	start, end := request.GetInterval(access.WithNow(s.Clk.Now()))
	req := types.CreateGrant{
		Id:       request.ID,
		Provider: accessRule.Target.ProviderID,
		With: types.CreateGrant_With{
			AdditionalProperties: make(map[string]string),
		},
		Subject: openapi_types.Email(q.Result.Email),
		Start:   iso8601.New(start),
		End:     iso8601.New(end),
	}

	//todo: rework this to be used safely around the codebase
	for k, v := range accessRule.Target.With {
		req.With.AdditionalProperties[k] = v
	}
	for k, v := range request.SelectedWith {
		req.With.AdditionalProperties[k] = v.Value
	}
	return req, nil
}
