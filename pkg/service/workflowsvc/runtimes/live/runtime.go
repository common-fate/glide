package live

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/sfn"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/common-fate/apikit/logger"
	ahTypes "github.com/common-fate/common-fate/accesshandler/pkg/types"
	"github.com/common-fate/common-fate/pkg/cfaws"
	"github.com/common-fate/common-fate/pkg/gevent"
	"github.com/common-fate/common-fate/pkg/targetgroupgranter"
)

type Runtime struct {
	StateMachineARN string
	AHClient        ahTypes.ClientWithResponsesInterface
	Eventbus        *gevent.Sender
}

func (r *Runtime) Grant(ctx context.Context, grant ahTypes.CreateGrant, isForTargetGroup bool) error {
	if isForTargetGroup {
		return r.grantTargetGroup(ctx, grant)
	}
	return r.grantProvider(ctx, grant)
}

func (r *Runtime) grantTargetGroup(ctx context.Context, grant ahTypes.CreateGrant) error {
	cfg, err := cfaws.ConfigFromContextOrDefault(ctx)
	if err != nil {
		return err
	}
	sfnClient := sfn.NewFromConfig(cfg)
	in := targetgroupgranter.WorkflowInput{Grant: grant}

	inJson, err := json.Marshal(in)
	if err != nil {
		return err
	}

	//running the step function
	sei := &sfn.StartExecutionInput{
		StateMachineArn: aws.String(r.StateMachineARN),
		Input:           aws.String(string(inJson)),
		Name:            &grant.Id,
	}

	//running the step function
	_, err = sfnClient.StartExecution(ctx, sei)
	return err

}

func (r *Runtime) grantProvider(ctx context.Context, grant ahTypes.CreateGrant) error {
	response, err := r.AHClient.PostGrantsWithResponse(ctx, grant)
	if err != nil {
		return err
	}
	if response.JSON201 != nil {
		return nil
	}
	if response.JSON400.Error != nil {
		return fmt.Errorf(*response.JSON400.Error)
	}
	logger.Get(ctx).Errorw("unhandled Access Handler response", "body", string(response.Body))
	return errors.New("unhandled response code from access provider service when granting")
}

func (r *Runtime) Revoke(ctx context.Context, grantID string, isForTargetGroup bool) error {
	if isForTargetGroup {
		return r.revokeTargetGroup(ctx, grantID)
	}
	return r.revokeProvider(ctx, grantID)
}

func (r *Runtime) revokeTargetGroup(ctx context.Context, grantID string) error {
	// @TODO
	return nil
}

func (r *Runtime) revokeProvider(ctx context.Context, grantID string) error {
	response, err := r.AHClient.PostGrantsRevokeWithResponse(ctx, grantID, ahTypes.PostGrantsRevokeJSONRequestBody{
		// @Note revoker ID is unused in the access handler code so it has been left empty here as it will not be included in this new runtime interface
		RevokerId: "",
	})
	if err != nil {
		return err
	}
	if response.JSON200 != nil {
		return nil
	}
	if response.JSON400.Error != nil {
		return fmt.Errorf(*response.JSON400.Error)
	}
	logger.Get(ctx).Errorw("unhandled Access Handler response", "body", string(response.Body))
	return errors.New("unhandled response code from access provider service when granting")
}
