package live

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	aws_config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	"go.uber.org/zap"

	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/cfaws"
	"github.com/common-fate/common-fate/pkg/gevent"
	"github.com/common-fate/common-fate/pkg/service/requestroutersvc"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/targetgroupgranter"
	"github.com/common-fate/ddb"
)

type GrantHandler interface {
	HandleRequest(ctx context.Context, in targetgroupgranter.InputEvent) (targetgroupgranter.GrantState, error)
}

type Runtime struct {
	StateMachineARN string
	Granter         GrantHandler

	Eventbus      *gevent.Sender
	DB            ddb.Storage
	RequestRouter *requestroutersvc.Service
}

func (r *Runtime) Grant(ctx context.Context, grant access.GroupTarget) error {

	cfg, err := cfaws.ConfigFromContextOrDefault(ctx)
	if err != nil {
		return err
	}
	sfnClient := sfn.NewFromConfig(cfg)
	in := targetgroupgranter.WorkflowInput{RequestAccessGroupTarget: grant}

	inJson, err := json.Marshal(in)
	if err != nil {
		return err
	}

	//running the step function
	sei := &sfn.StartExecutionInput{
		StateMachineArn: aws.String(r.StateMachineARN),
		Input:           aws.String(string(inJson)),
		Name:            &grant.ID,
	}

	//running the step function
	_, err = sfnClient.StartExecution(ctx, sei)
	return err

}
func (r *Runtime) Revoke(ctx context.Context, grantID string) error {
	// we can grab all this from the execution input for the step function we will use this as the source of truth
	c, err := aws_config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}
	sfnClient := sfn.NewFromConfig(c)

	// build the execution ARN
	exeARN := BuildExecutionARN(r.StateMachineARN, grantID)

	out, err := sfnClient.DescribeExecution(ctx, &sfn.DescribeExecutionInput{ExecutionArn: aws.String(exeARN)})
	if err != nil {
		return err
	}

	// build the previous grant from the execution input
	var input targetgroupgranter.WorkflowInput

	err = json.Unmarshal([]byte(*out.Input), &input)
	if err != nil {
		return err
	}

	// if the state function is in the active state then we will stop the execution
	statefn, err := sfnClient.GetExecutionHistory(ctx, &sfn.GetExecutionHistoryInput{ExecutionArn: &exeARN})
	if err != nil {
		return err
	}
	tgq := storage.GetTargetGroup{
		ID: input.RequestAccessGroupTarget.TargetGroupID,
	}

	_, err = r.DB.Query(ctx, &tgq)
	if err != nil {
		return err
	}

	lastState := statefn.Events[len(statefn.Events)-1]

	// if the state of the grant is in the active state
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

		zap.S().Infow("calling invoke lambda to revoke")

		_, err = r.Granter.HandleRequest(ctx, targetgroupgranter.InputEvent{
			Action:                   targetgroupgranter.DEACTIVATE,
			RequestAccessGroupTarget: input.RequestAccessGroupTarget,
			State:                    gs.State,
		})
		if err != nil {
			zap.S().Errorw("failed to deactivate grant", "err", err)
			return err
		}
	}

	_, err = sfnClient.StopExecution(ctx, &sfn.StopExecutionInput{ExecutionArn: &exeARN})
	if err != nil {
		return err
	}

	return nil

}

func BuildExecutionARN(stateMachineARN string, grantID string) string {

	splitARN := strings.Split(stateMachineARN, ":")

	//position 5 is the location of the arn type
	splitARN[5] = "execution"
	splitARN = append(splitARN, grantID)

	return strings.Join(splitARN, ":")

}
