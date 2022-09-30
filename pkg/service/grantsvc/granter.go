package grantsvc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/benbjohnson/clock"
	"github.com/segmentio/ksuid"

	"github.com/common-fate/apikit/logger"

	"github.com/common-fate/ddb"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providerregistry"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers"
	ahTypes "github.com/common-fate/granted-approvals/accesshandler/pkg/types"
	"github.com/common-fate/granted-approvals/pkg/access"
	"github.com/common-fate/granted-approvals/pkg/deploy"
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
	AHClient           ahTypes.ClientWithResponsesInterface
	DB                 ddb.Storage
	Clock              clock.Clock
	EventBus           *gevent.Sender
	accessTokenChecker accessTokenChecker
}

type GranterOpts struct {
	AHClient         ahTypes.ClientWithResponsesInterface
	DB               ddb.Storage
	Clock            clock.Clock
	EventBus         *gevent.Sender
	DeploymentConfig deploy.DeployConfigReader
}

// New creates a new Granter service.
func New(opts GranterOpts) *Granter {
	return &Granter{
		AHClient: opts.AHClient,
		DB:       opts.DB,
		Clock:    opts.Clock,
		EventBus: opts.EventBus,
		accessTokenChecker: registryAccessTokenChecker{
			DeploymentConfig: opts.DeploymentConfig,
			Registry:         providerregistry.Registry(),
		},
	}
}

type CreateGrantOpts struct {
	Request    access.Request
	AccessRule rule.AccessRule
}

type RevokeGrantOpts struct {
	Request   access.Request
	RevokerID string
}

func (g *Granter) RevokeGrant(ctx context.Context, opts RevokeGrantOpts) (*access.Request, error) {
	if opts.Request.Grant == nil {
		return nil, ErrNoGrant
	}
	//Cannot request to revoke/cancel grant if it is not active or pending (state function has been created and executed)
	canRevoke := opts.Request.Grant.Status == ahTypes.GrantStatusACTIVE || opts.Request.Grant.Status == ahTypes.GrantStatusPENDING

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
		opts.Request.Grant.Status = ahTypes.GrantStatusREVOKED
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

// validate grant runs all the checks that will need to occur when creating a real grant to validate its success
func (g *Granter) ValidateGrant(ctx context.Context, opts CreateGrantOpts) ([]ahTypes.GrantValidation, error) {
	q := &storage.GetUser{
		ID: opts.Request.RequestedBy,
	}
	_, err := g.DB.Query(ctx, q)
	if err != nil {
		return []ahTypes.GrantValidation{}, err
	}
	start, end := opts.Request.GetInterval(access.WithNow(g.Clock.Now()))

	req := ahTypes.ValidateRequestToProviderJSONRequestBody{
		Id:       opts.Request.ID,
		Provider: opts.AccessRule.Target.ProviderID,
		With: ahTypes.CreateGrant_With{
			AdditionalProperties: make(map[string]string),
		},
		Subject: openapi_types.Email(q.Result.Email),
		Start:   iso8601.New(start),
		End:     iso8601.New(end),
	}
	for k, v := range opts.AccessRule.Target.With {
		req.With.AdditionalProperties[k] = v
	}
	for k, v := range opts.Request.SelectedWith {
		req.With.AdditionalProperties[k] = v.Value
	}

	res, err := g.AHClient.ValidateRequestToProviderWithResponse(ctx, req)
	if err != nil {
		return []ahTypes.GrantValidation{}, err
	}

	if res.JSON200 == nil {
		//there was an error, handle it
		return []ahTypes.GrantValidation{}, fmt.Errorf("error validating grant: %s", res.Body)
	} else {
		var respBody ahTypes.GrantValidationResponse
		err := json.Unmarshal(res.Body, &respBody)
		if err != nil {
			return []ahTypes.GrantValidation{}, err
		}

		//check if any failed states
		validationFailed := false
		for _, v := range respBody.Validation {
			if v.Status == "ERROR" {
				validationFailed = true
			}
		}

		if validationFailed {
			return respBody.Validation, GrantValidationError{}

		}
		return respBody.Validation, nil

	}

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

	// check whether the Access Provider requires an Access Token to be generated - we'll create one if it does.
	// check now before we actually provision the access, so that we can return early if we fail.
	requiresAccessToken, err := g.accessTokenChecker.NeedsAccessToken(ctx, opts.AccessRule.Target.ProviderID)
	if err != nil {
		return nil, err
	}

	start, end := opts.Request.GetInterval(access.WithNow(g.Clock.Now()))
	req := ahTypes.PostGrantsJSONRequestBody{
		Id:       opts.Request.ID,
		Provider: opts.AccessRule.Target.ProviderID,
		With: ahTypes.CreateGrant_With{
			AdditionalProperties: make(map[string]string),
		},
		Subject: openapi_types.Email(q.Result.Email),
		Start:   iso8601.New(start),
		End:     iso8601.New(end),
	}
	for k, v := range opts.AccessRule.Target.With {
		req.With.AdditionalProperties[k] = v
	}
	for k, v := range opts.Request.SelectedWith {
		req.With.AdditionalProperties[k] = v.Value
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

		if requiresAccessToken {
			logger.Get(ctx).Infow("creating access token for request", "request.id", opts.Request.ID)
			at := access.AccessToken{
				RequestID: opts.Request.ID,
				Token:     ksuid.New().String(),
				Start:     res.JSON201.Grant.Start.Time,
				End:       res.JSON201.Grant.End.Time,
				CreatedAt: now,
			}
			err = g.DB.Put(ctx, &at)
			if err != nil {
				return nil, err
			}
		}

		return &opts.Request, nil
	}

	if res.JSON400.Error != nil {
		return nil, fmt.Errorf(*res.JSON400.Error)
	}
	logger.Get(ctx).Errorw("unhandled Access Handler response", "body", string(res.Body))
	return nil, errors.New("unhandled response code")
}

// accessTokenCheckers check whether a provider needs an access token generated.
type accessTokenChecker interface {
	NeedsAccessToken(ctx context.Context, providerID string) (bool, error)
}

type registryAccessTokenChecker struct {
	DeploymentConfig deploy.DeployConfigReader
	Registry         providerregistry.ProviderRegistry
}

// providerRequiresAccessToken looks up the provider in our registry.
// If the provider implements RequiresAccessToken() and it's true, this function returns true.
// Otherwise, it returns false.
// Returns an error if we can't look up the provider.
func (r registryAccessTokenChecker) NeedsAccessToken(ctx context.Context, providerID string) (bool, error) {
	pm, err := r.DeploymentConfig.ReadProviders(ctx)
	if err != nil {
		return false, err
	}
	provider, ok := pm[providerID]
	if !ok {
		return false, fmt.Errorf("could not find provider %s in deployment config", providerID)
	}
	p, err := r.Registry.LookupByUses(provider.Uses)
	if err != nil {
		return false, err
	}
	if at, ok := p.Provider.(providers.AccessTokener); ok && at.RequiresAccessToken() {
		return true, nil
	}
	return false, nil
}
