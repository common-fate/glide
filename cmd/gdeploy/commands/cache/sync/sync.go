package sync

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/briandowns/spinner"
	"github.com/common-fate/clio"
	"github.com/common-fate/clio/clierr"
	"github.com/common-fate/common-fate/pkg/cfaws"
	"github.com/common-fate/common-fate/pkg/deploy"
	"github.com/urfave/cli/v2"
)

var SyncCommand = cli.Command{
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

		if o.CacheSyncFunctionName == "" {
			return clierr.New("The sync function name is not yet available. You may need to update your deployment to use this feature.")
		}
		si := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		si.Suffix = " invoking cache sync lambda function"
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
		b, err := json.Marshal(res)
		if err != nil {
			return err
		}
		clio.Debugf("cache sync lamda invoke response: %s", string(b))
		if res.FunctionError != nil {
			return fmt.Errorf("cache sync failed with lambda execution error: %s", *res.FunctionError)
		} else if res.StatusCode == 200 {

			clio.Successf("Successfully synced the cache")
		} else {
			return fmt.Errorf("cache sync failed with lambda invoke status code: %d", res.StatusCode)
		}
		return nil
	}}
