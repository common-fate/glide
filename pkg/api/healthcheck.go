package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	lambdaTypes "github.com/aws/aws-sdk-go-v2/service/lambda/types"

	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/briandowns/spinner"
	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/clio"
	"github.com/common-fate/common-fate/pkg/cfaws"
	"github.com/common-fate/common-fate/pkg/deploy"
	"github.com/common-fate/common-fate/pkg/storage"
)

func (a *API) AdminRunHealthcheck(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	dc, err := deploy.LoadConfig(deploy.DefaultFilename)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	o, err := dc.LoadOutput(ctx)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	cfg, err := cfaws.ConfigFromContextOrDefault(ctx)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	// if o.HealthcheckFunctionName == "" {
	// 	return clierr.New("The healthcheck function name is not yet available. You may need to update your deployment to use this feature.")
	// }

	si := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	si.Suffix = " invoking healthcheck sync lambda function"
	si.Writer = os.Stderr
	si.Start()

	lambdaClient := lambda.NewFromConfig(cfg)
	res, err := lambdaClient.Invoke(ctx, &lambda.InvokeInput{
		FunctionName:   &o.HealthcheckFunctionName,
		InvocationType: lambdaTypes.InvocationTypeRequestResponse,
		Payload:        []byte("{}"),
	})
	si.Stop()
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	b, err := json.Marshal(res)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	clio.Debugf("healthcheck sync lambda invoke response: %s", string(b))
	if res.FunctionError != nil {
		apio.Error(ctx, w, fmt.Errorf("healthcheck sync failed with lambda execution error: %s", *res.FunctionError))
		return
	} else if res.StatusCode == 200 {

		clio.Successf("Successfully synced the healthcheck")
	} else {
		apio.Error(ctx, w, fmt.Errorf("healthcheck sync failed with lambda invoke status code: %d", res.StatusCode))
		return
	}

	// now run a fetch
	listHandlers := storage.ListHandlers{}

	_, err = a.DB.Query(ctx, &listHandlers)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	healthyCount := 0
	unhealthyCount := 0

	for _, deployment := range listHandlers.Result {
		if deployment.Healthy {
			healthyCount++
		} else {
			unhealthyCount++
		}
	}

	clio.Log("healthcheck result")
	clio.Logf("healthy: %d", healthyCount)
	clio.Logf("unhealthy: %d", unhealthyCount)

	apio.JSON(ctx, w, nil, http.StatusNoContent)
}
