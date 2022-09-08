package lambda

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sfn"

	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/types"
	"github.com/common-fate/granted-approvals/pkg/gevent"
)

// WorkflowInput is the input to the Step Functions workflow execution
type WorkflowInput struct {
	Grant types.Grant `json:"grant"`
}

// CreateGrant creates a new grant.
func (r *Runtime) CreateGrant(ctx context.Context, vcg types.ValidCreateGrant) (types.Grant, error) {
	grant := types.NewGrant(vcg)
	logger.Get(ctx).Infow("creating grant", "grant", grant)

	//setting the input for the step function lambda

	in := WorkflowInput{Grant: grant}

	logger.Get(ctx).Infow("constructed workflow input", "input", in)

	inJson, err := json.Marshal(in)
	if err != nil {
		return types.Grant{}, err
	}
	c, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return types.Grant{}, err
	}

	sfnClient := sfn.NewFromConfig(c)

	//running the step function
	sei := &sfn.StartExecutionInput{
		StateMachineArn: aws.String(r.StateMachineARN),
		Input:           aws.String(string(inJson)),
		Name:            &grant.ID,
	}

	//running the step function
	_, err = sfnClient.StartExecution(ctx, sei)
	if err != nil {
		return types.Grant{}, err
	}
	eventsBus, err := gevent.NewSender(ctx, gevent.SenderOpts{EventBusARN: r.EventBusArn})
	if err != nil {
		return types.Grant{}, err
	}
	evt := &gevent.GrantCreated{Grant: grant}
	err = eventsBus.Put(ctx, evt)
	if err != nil {
		return types.Grant{}, err
	}

	return grant, nil
}
