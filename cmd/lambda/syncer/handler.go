package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/common-fate/pkg/config"
	"github.com/common-fate/common-fate/pkg/deploy"
	"github.com/common-fate/common-fate/pkg/identity/identitysync"
	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
	"go.uber.org/zap"
)

func main() {
	var cfg config.SyncConfig
	ctx := context.Background()
	_ = godotenv.Load()

	err := envconfig.Process(ctx, &cfg)
	if err != nil {
		panic(err)
	}

	ic, err := deploy.UnmarshalFeatureMap(cfg.IdentitySettings)
	if err != nil {
		panic(err)
	}
	//set up the sync handler
	syncer, err := identitysync.NewIdentitySyncer(ctx, identitysync.SyncOpts{
		TableName:      cfg.TableName,
		IdpType:        cfg.IdpProvider,
		UserPoolId:     cfg.UserPoolId,
		IdentityConfig: ic,
	})
	if err != nil {
		panic(err)
	}
	log, err := logger.Build(cfg.LogLevel)
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(log.Desugar())
	zap.S().Infow("starting sync", "config", ic, "idp.type", cfg.IdpProvider)
	lambda.Start(syncer.Sync)
}
