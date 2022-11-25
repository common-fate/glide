package main

import (
	"context"

	"go.uber.org/zap"

	"github.com/aws/aws-lambda-go/cfn"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/common-fate/pkg/config"
	"github.com/common-fate/common-fate/pkg/deploy"
	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

func main() {
	lambda.Start(cfn.LambdaWrap(func(ctx context.Context, e cfn.Event) (physicalResourceID string, data map[string]interface{}, err error) {
		physicalResourceID = "grantedFrontendConfigurer"
		var cfg config.FrontendDeployerConfig
		_ = godotenv.Load()
		err = envconfig.Process(ctx, &cfg)
		if err != nil {
			return
		}
		gLog, err := logger.Build(cfg.LogLevel)
		if err != nil {
			panic(err)
		}
		zap.ReplaceGlobals(gLog.Desugar())

		// Log the event - so that if we get stuck in the DELETE_FAILED state
		// we can manually trigger a deletion of the resource.
		// More info: https://aws.amazon.com/premiumsupport/knowledge-center/cloudformation-lambda-resource-delete/
		log := zap.S().With("config", cfg, "event", e)
		log.Info("running frontend configurer")
		// publish the frontend for every update and delete
		if e.RequestType == cfn.RequestCreate || e.RequestType == cfn.RequestUpdate {
			// attempt to publish the aws exports for frontend configuration
			// This config unpacking is annoying but it has to happen somewhere
			err = deploy.DeployProductionFrontend(ctx, cfg)
		}
		return
	}))

}
