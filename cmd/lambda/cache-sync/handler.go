package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/common-fate/internal"
	"github.com/common-fate/common-fate/pkg/cachesync"
	"github.com/common-fate/common-fate/pkg/config"
	"github.com/common-fate/common-fate/pkg/service/cachesvc"
	"github.com/common-fate/common-fate/pkg/service/requestroutersvc"
	"github.com/common-fate/ddb"
	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
	"go.uber.org/zap"
)

func main() {
	var cfg config.CacheSyncConfig
	ctx := context.Background()
	_ = godotenv.Load()

	err := envconfig.Process(ctx, &cfg)
	if err != nil {
		panic(err)
	}
	db, err := ddb.New(ctx, cfg.TableName)
	if err != nil {
		panic(err)
	}
	ahc, err := internal.BuildAccessHandlerClient(ctx, internal.BuildAccessHandlerClientOpts{Region: cfg.Region, AccessHandlerURL: cfg.AccessHandlerURL})
	if err != nil {
		panic(err)
	}

	syncer := cachesync.CacheSyncer{
		DB:                  db,
		AccessHandlerClient: ahc,
		Cache: cachesvc.Service{
			DB:                  db,
			AccessHandlerClient: ahc,
			RequestRouter: &requestroutersvc.Service{
				DB: db,
			},
		},
	}
	log, err := logger.Build(cfg.LogLevel)
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(log.Desugar())
	zap.S().Infow("starting cache sync", "config", cfg)
	lambda.Start(syncer.Sync)
}
