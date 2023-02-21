package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/common-fate/pkg/config"
	"github.com/common-fate/common-fate/pkg/service/healthchecksvc"
	"github.com/common-fate/ddb"
	"github.com/sethvargo/go-envconfig"
	"go.uber.org/zap"
)

func main() {
	var cfg config.HealthCheckerConfig
	ctx := context.Background()

	err := envconfig.Process(ctx, &cfg)
	if err != nil {
		panic(err)
	}
	db, err := ddb.New(ctx, cfg.TableName)
	if err != nil {
		panic(err)
	}

	healthchecker := healthchecksvc.Service{
		DB: db,
	}
	log, err := logger.Build(cfg.LogLevel)
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(log.Desugar())
	zap.S().Infow("starting healthchecker check", "config", cfg)
	lambda.Start(healthchecker.Check)
}
