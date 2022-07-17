package users

import (
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/briandowns/spinner"
	"github.com/common-fate/granted-approvals/pkg/cfaws"
	"github.com/common-fate/granted-approvals/pkg/clio"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/urfave/cli/v2"
)

var syncCommand = cli.Command{
	Name: "sync",
	Action: func(c *cli.Context) error {
		ctx := c.Context

		dc, err := deploy.ConfigFromContext(ctx)
		if err != nil {
			return err
		}
		o, err := dc.LoadOutput(ctx)
		if err != nil {
			return err
		}

		cfg, err := cfaws.ConfigFromContextOrDefault(ctx)
		if err != nil {
			return err
		}

		if o.IdpSyncFunctionName == "" {
			return clio.NewCLIError("The sync function name is not yet available. You may need to update your deployment to use this feature.")
		}
		si := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		si.Suffix = " invoking IDP sync lambda function"
		si.Writer = os.Stderr
		si.Start()

		lambdaClient := lambda.NewFromConfig(cfg)
		res, err := lambdaClient.Invoke(ctx, &lambda.InvokeInput{
			FunctionName:   &o.IdpSyncFunctionName,
			InvocationType: types.InvocationTypeRequestResponse,
			Payload:        []byte("{}"),
		})
		si.Stop()
		if err != nil {
			return err
		}

		clio.Info("Lambda execution completed with status: %d. ", res.StatusCode)
		if res.FunctionError != nil {
			clio.Error("Lambda returned execution error: %s", *res.FunctionError)
		}
		return nil
	}}
