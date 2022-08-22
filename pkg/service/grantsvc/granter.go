package grantsvc

import (
	"context"
	"errors"
	"fmt"

	"github.com/benbjohnson/clock"

	"github.com/common-fate/apikit/logger"

	"github.com/common-fate/ddb"
	ahTypes "github.com/common-fate/granted-approvals/accesshandler/pkg/types"
	"github.com/common-fate/granted-approvals/pkg/access"
	"github.com/common-fate/granted-approvals/pkg/gevent"
	"github.com/common-fate/granted-approvals/pkg/rule"
	"github.com/common-fate/granted-approvals/pkg/storage"
	"github.com/common-fate/granted-approvals/pkg/storage/dbupdate"
	"github.com/common-fate/granted-approvals/pkg/types"
	"github.com/common-fate/iso8601"
	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
)

type UserGetter interface {
	GetUserBySub(ctx context.Context, sub string) (*types.User, error)
}

//go:generate go run github.com/golang/mock/mockgen -destination=mocks/mock_accesshandler_client.go -package=mocks . AHClient
type AHClient interface {
	ahTypes.ClientWithResponsesInterface
}

// Granter has logic to integrate with the Access Handler.
type Granter struct {
	AHClient ahTypes.ClientWithResponsesInterface
	DB       ddb.Storage
	Clock    clock.Clock
	EventBus *gevent.Sender
}

type CreateGrantOpts struct {
	Request    access.Request
	AccessRule rule.AccessRule
}

type RevokeGrantOpts struct {
	Request   access.Request
	RevokerID string
}

// NewGranter creates a new Granter instance
func NewGranter(client ahTypes.ClientWithResponsesInterface, db ddb.Storage, clock clock.Clock, eventBus *gevent.Sender) (*Granter, error) {

	return &Granter{client, db, clock, eventBus}, nil
}

func (g *Granter) RevokeGrant(ctx context.Context, opts RevokeGrantOpts) (*access.Request, error) {
	if opts.Request.Grant == nil {
		return nil, ErrNoGrant
	}
	//Cannot request to revoke/cancel grant if it is not active or pending (state function has been created and executed)
	canRevoke := opts.Request.Grant.Status == ahTypes.ACTIVE || opts.Request.Grant.Status == ahTypes.PENDING

	if !canRevoke || opts.Request.Grant.End.Before(g.Clock.Now()) {
		return nil, ErrGrantInactive
	}
	res, err := g.AHClient.PostGrantsRevokeWithResponse(ctx, opts.Request.ID, ahTypes.PostGrantsRevokeJSONRequestBody{
		RevokerId: opts.RevokerID,
	})
	if err != nil {
		return nil, err
	}

	if res.JSON200 != nil {
		oldStatus := opts.Request.Grant.Status
		opts.Request.Grant.Status = ahTypes.REVOKED
		opts.Request.Grant.UpdatedAt = g.Clock.Now()
		items, err := dbupdate.GetUpdateRequestItems(ctx, g.DB, opts.Request)
		if err != nil {
			return nil, err
		}

		//create a request event for audit loggging request change
		requestEvent := access.NewGrantStatusChangeEvent(opts.Request.ID, opts.Request.Grant.UpdatedAt, &opts.RevokerID, oldStatus, opts.Request.Grant.Status)

		items = append(items, &requestEvent)

		err = g.DB.PutBatch(ctx, items...)
		if err != nil {
			return nil, err
		}

		// Emit an event for the grant revoke
		// We have chosen to emit events from the approvals app for grant revocation rather than from the access handler because we are using a syncronous API.
		// All effects from revoking will be implemented in this syncronous api rather than triggered from the events.
		// So we update the grant status here and save the grant before emitting the event
		err = g.EventBus.Put(ctx, gevent.GrantRevoked{Grant: opts.Request.Grant.ToAHGrant(opts.Request.ID)})
		if err != nil {
			return nil, err
		}
		return &opts.Request, nil
	}

	if res.JSON400 != nil {
		logger.Get(ctx).Errorw("Invalid request", "body", string(res.Body))

		return nil, fmt.Errorf(*res.JSON400.Error)
	}

	if res.JSON500 != nil {
		logger.Get(ctx).Errorw("Internal server error", "body", string(res.Body))

		return nil, fmt.Errorf(*res.JSON500.Error)
	}
	logger.Get(ctx).Errorw("unhandled Access Handler response", "body", string(res.Body))
	return nil, errors.New("unhandled response code")

}

// CreateGrant creates a Grant in the Access Handler, it does not update the approvals app database.
// the returned Request will contain the newly created grant
func (g *Granter) CreateGrant(ctx context.Context, opts CreateGrantOpts) (*access.Request, error) {

	q := &storage.GetUser{
		ID: opts.Request.RequestedBy,
	}
	_, err := g.DB.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	start, end := opts.Request.GetInterval(access.WithNow(g.Clock.Now()))
	req := ahTypes.PostGrantsJSONRequestBody{
		Id:       opts.Request.ID,
		Provider: opts.AccessRule.Target.ProviderID,
		With: ahTypes.CreateGrant_With{
			AdditionalProperties: opts.AccessRule.Target.With,
		},
		Subject: openapi_types.Email(q.Result.Email),
		Start:   iso8601.New(start),
		End:     iso8601.New(end),
	}

	res, err := g.AHClient.PostGrantsWithResponse(ctx, req)
	if err != nil {
		return nil, err
	}

	// on success we create a grant item in dynamo db
	if res.JSON201 != nil {
		now := g.Clock.Now()
		opts.Request.Grant = &access.Grant{
			Provider:  res.JSON201.Grant.Provider,
			Subject:   string(res.JSON201.Grant.Subject),
			Start:     res.JSON201.Grant.Start.Time,
			End:       res.JSON201.Grant.End.Time,
			Status:    res.JSON201.Grant.Status,
			With:      res.JSON201.Grant.With,
			CreatedAt: now,
			UpdatedAt: now,
		}

		//save access token to dynamo
		newAccessToken := access.AccessToken{RequestId: opts.Request.ID, Token: *res.JSON201.AdditionalProperties.AccessToken}

		err = g.DB.Put(ctx, &newAccessToken)
		if err != nil {
			return nil, err
		}

		return &opts.Request, nil
	}

	if res.JSON400.Error != nil {
		return nil, fmt.Errorf(*res.JSON400.Error)
	}
	logger.Get(ctx).Errorw("unhandled Access Handler response", "body", string(res.Body))
	return nil, errors.New("unhandled response code")
}
