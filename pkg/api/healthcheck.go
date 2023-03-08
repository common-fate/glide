package api

import (
	"context"
	"fmt"
	"net/http"

	lambdaTypes "github.com/aws/aws-sdk-go-v2/service/lambda/types"

	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/common-fate/pkg/cfaws"
	"github.com/common-fate/common-fate/pkg/deploy"
	"github.com/common-fate/common-fate/pkg/storage"
)

func (a *API) AdminRunHealthcheck(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	log := logger.Get(ctx)
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

	lambdaClient := lambda.NewFromConfig(cfg)
	res, err := lambdaClient.Invoke(ctx, &lambda.InvokeInput{
		FunctionName:   &o.HealthcheckFunctionName,
		InvocationType: lambdaTypes.InvocationTypeRequestResponse,
		Payload:        []byte("{}"),
	})
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	if res.FunctionError != nil {
		apio.Error(ctx, w, fmt.Errorf("healthcheck sync failed with lambda execution error: %s", *res.FunctionError))
		return
	} else if res.StatusCode == 200 {

		log.Info("Successfully synced the healthcheck")
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

	apio.JSON(ctx, w, nil, http.StatusNoContent)
}
