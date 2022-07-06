package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/granted-approvals/pkg/config"
	"github.com/common-fate/granted-approvals/pkg/eventhandler"

	"github.com/common-fate/ddb"
	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
	"go.uber.org/zap"
)

func main() {
	var cfg config.EventHandlerConfig
	ctx := context.Background()
	_ = godotenv.Load()
	err := envconfig.Process(ctx, &cfg)
	if err != nil {
		panic(err)
	}
	db, err := ddb.New(ctx, cfg.DynamoTable)
	if err != nil {
		panic(err)
	}
	eventHandler, err := eventhandler.New(ctx, db)
	if err != nil {
		panic(err)
	}
	log, err := logger.Build(cfg.LogLevel)
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(log.Desugar())
	zap.S().Infow("starting event handler with configuration", "config", cfg)
	lambda.Start(eventHandler.HandleEvent)
}
