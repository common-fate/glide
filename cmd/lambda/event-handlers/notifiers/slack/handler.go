package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/ddb"
	"github.com/common-fate/granted-approvals/pkg/config"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	slacknotifier "github.com/common-fate/granted-approvals/pkg/notifiers/slack"
	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
	"go.uber.org/zap"
)

func main() {
	var cfg config.SlackNotifierConfig
	ctx := context.Background()
	_ = godotenv.Load()

	err := envconfig.Process(ctx, &cfg)
	if err != nil {
		panic(err)
	}
	log, err := logger.Build(cfg.LogLevel)
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(log.Desugar())
	db, err := ddb.New(ctx, cfg.DynamoTable)
	if err != nil {
		panic(err)
	}

	s, err := deploy.UnmarshalSlack(cfg.SlackSettings)
	if err != nil {
		panic(err)
	}

	zap.S().Infow("starting notifier with configuration", "config", cfg)
	err = config.LoadAndReplaceSSMValues(ctx, &s)
	if err != nil {
		panic(err)
	}
	notifier := &slacknotifier.Notifier{
		DB:          db,
		FrontendURL: cfg.FrontendURL,
		SlackConfig: s,
	}

	lambda.Start(notifier.HandleEvent)
}
