package grantsvcv2

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/service/sfn"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/benbjohnson/clock"

	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/common-fate/pkg/cfaws"

	ahTypes "github.com/common-fate/common-fate/accesshandler/pkg/types"
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/gevent"
	"github.com/common-fate/common-fate/pkg/pdk"
	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/targetgroupgranter"
	"github.com/common-fate/ddb"
	"github.com/common-fate/iso8601"
	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
)

// Granter has logic to integrate with the Access Handler.
type Granter struct {
	db              ddb.Storage
	clock           clock.Clock
	eventBus        *gevent.Sender
	stateMachineARN string
	client          *sfn.Client
}

type GranterOpts struct {
	DB              ddb.Storage
	Clock           clock.Clock
	EventBus        *gevent.Sender
	StateMachineARN string
}

// New creates a new Granter service.
func New(ctx context.Context, opts GranterOpts) (*Granter, error) {
	cfg, err := cfaws.ConfigFromContextOrDefault(ctx)
	if err != nil {
		return nil, err
	}
	sfnClient := sfn.NewFromConfig(cfg)
	g := Granter{
		db:              opts.DB,
		clock:           opts.Clock,
		eventBus:        opts.EventBus,
		client:          sfnClient,
		stateMachineARN: opts.StateMachineARN,
	}
	return &g, nil
}

type CreateGrantOpts struct {
	Request    access.Request
	AccessRule rule.AccessRule
}

type RevokeGrantOpts struct {
	Request   access.Request
	RevokerID string
}

type WorkflowInput struct {
	Grant targetgroupgranter.Grant `json:"grant"`
}

// CreateGrant creates a Grant in the Access Handler, it does not update the Common Fate app database.
// the returned Request will contain the newly created grant
func (g *Granter) CreateGrant(ctx context.Context, opts CreateGrantOpts) (*access.Request, error) {
	grant, err := g.prepareCreateGrantRequest(ctx, opts)
	if err != nil {
		return nil, err
	}

	logger.Get(ctx).Infow("creating grant", "grant", grant)

	//setting the input for the step function lambda

	in := WorkflowInput{Grant: grant}

	logger.Get(ctx).Infow("constructed workflow input", "input", in)

	inJson, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}

	//running the step function
	sei := &sfn.StartExecutionInput{
		StateMachineArn: aws.String(g.stateMachineARN),
		Input:           aws.String(string(inJson)),
		Name:            &grant.ID,
	}

	//running the step function
	_, err = g.client.StartExecution(ctx, sei)
	if err != nil {
		return nil, err
	}
	now := g.clock.Now()

	// @TODO the v1 grant type is baked in here a little, I have adapted it for v2
	opts.Request.Grant = &access.Grant{
		Provider:  grant.TargetGroup,
		Subject:   string(grant.Subject),
		Start:     grant.Start.Time,
		End:       grant.End.Time,
		Status:    ahTypes.GrantStatusPENDING,
		With:      ahTypes.Grant_With{AdditionalProperties: grant.Target.Arguments},
		CreatedAt: now,
		UpdatedAt: now,
	}

	return &opts.Request, nil
}

// prepareCreateGrantRequest converts opts into a CreateGrant struct for access handler requests
func (g *Granter) prepareCreateGrantRequest(ctx context.Context, opts CreateGrantOpts) (targetgroupgranter.Grant, error) {
	q := &storage.GetUser{
		ID: opts.Request.RequestedBy,
	}
	_, err := g.db.Query(ctx, q)
	if err != nil {
		return targetgroupgranter.Grant{}, err
	}
	target := make(map[string]string)
	start, end := opts.Request.GetInterval(access.WithNow(g.clock.Now()))
	req := targetgroupgranter.Grant{
		ID:          opts.Request.ID,
		TargetGroup: opts.AccessRule.Target.TargetGroupID,
		Subject:     openapi_types.Email(q.Result.Email),
		Start:       iso8601.New(start),
		End:         iso8601.New(end),
	}

	for k, v := range opts.AccessRule.Target.With {
		target[k] = v
	}
	for k, v := range opts.Request.SelectedWith {
		target[k] = v.Value
	}
	req.Target = pdk.NewDefaultModeTarget(target)
	return req, nil
}
