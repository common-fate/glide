package live

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	aws_config "github.com/aws/aws-sdk-go-v2/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	"github.com/common-fate/apikit/logger"
	ahTypes "github.com/common-fate/common-fate/accesshandler/pkg/types"
	"github.com/common-fate/common-fate/pkg/cfaws"
	"github.com/common-fate/common-fate/pkg/gevent"
	"github.com/common-fate/common-fate/pkg/handler"
	"github.com/common-fate/common-fate/pkg/service/requestroutersvc"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/targetgroupgranter"
	"github.com/common-fate/ddb"
	"github.com/common-fate/provider-registry-sdk-go/pkg/msg"
)

type Runtime struct {
	StateMachineARN string
	AHClient        ahTypes.ClientWithResponsesInterface
	Eventbus        *gevent.Sender
	DB              ddb.Storage
	RequestRouter   *requestroutersvc.Service
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

func BuildExecutionARN(stateMachineARN string, grantID string) string {

	splitARN := strings.Split(stateMachineARN, ":")

	//position 5 is the location of the arn type
	splitARN[5] = "execution"
	splitARN = append(splitARN, grantID)

	return strings.Join(splitARN, ":")

}

func (r *Runtime) revokeTargetGroup(ctx context.Context, grantID string) error {
	//we can grab all this from the execution input for the step function we will use this as the source of truth
	c, err := aws_config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}
	sfnClient := sfn.NewFromConfig(c)

	//build the execution ARN
	exeARN := BuildExecutionARN(r.StateMachineARN, grantID)

	out, err := sfnClient.DescribeExecution(ctx, &sfn.DescribeExecutionInput{ExecutionArn: aws.String(exeARN)})
	if err != nil {
		return err
	}

	//build the previous grant from the execution input
	var input targetgroupgranter.WorkflowInput

	err = json.Unmarshal([]byte(*out.Input), &input)
	if err != nil {
		return err
	}

	//if the state function is in the active state then we will stop the execution
	statefn, err := sfnClient.GetExecutionHistory(ctx, &sfn.GetExecutionHistoryInput{ExecutionArn: &exeARN})
	if err != nil {
		return err
	}
	tgq := storage.GetTargetGroup{
		ID: input.Grant.Provider,
	}

	_, err = r.DB.Query(ctx, &tgq)
	if err != nil {
		return err
	}
	routeResult, err := r.RequestRouter.Route(ctx, *tgq.Result)
	if err != nil {
		return err
	}
	runtime, err := handler.GetRuntime(ctx, routeResult.Handler)
	if err != nil {
		return err
	}
	lastState := statefn.Events[len(statefn.Events)-1]

	//if the state of the grant is in the active state
	if lastState.Type == "WaitStateEntered" && *lastState.StateEnteredEventDetails.Name == "Wait for Window End" {

		// Pull the state from the output of the activate step so it can be used when revoking access
		exitActivateStepEvent := statefn.Events[len(statefn.Events)-2]
		if exitActivateStepEvent.Type != "TaskStateExited" || exitActivateStepEvent.StateExitedEventDetails == nil {
			return errors.New("unexpected workflow state")
		}

		var gs targetgroupgranter.GrantState
		err = json.Unmarshal([]byte(aws.ToString(exitActivateStepEvent.StateExitedEventDetails.Output)), &gs)
		if err != nil {
			return err
		}
		//call the provider revoke
		req := msg.Revoke{
			Subject: string(input.Grant.Subject),
			Target: msg.Target{
				Kind:      routeResult.Route.Kind,
				Arguments: input.Grant.With.AdditionalProperties,
			},
			Request: msg.AccessRequest{
				ID: grantID,
			},
			State: gs.State,
		}

		err = runtime.Revoke(ctx, req)
		if err != nil {
			return err
		}
	}

	_, err = sfnClient.StopExecution(ctx, &sfn.StopExecutionInput{ExecutionArn: &exeARN})
	if err != nil {
		return err
	}

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

	switch response.StatusCode() {
	case http.StatusOK:
		return nil
	case http.StatusBadRequest:
		return fmt.Errorf(*response.JSON400.Error)
	case http.StatusInternalServerError:
		return fmt.Errorf(*response.JSON500.Error)
	default:
		logger.Get(ctx).Errorw("unhandled Access Handler response", "body", string(response.Body))
		return errors.New("unhandled response code from access provider service when revoking access")
	}

}
