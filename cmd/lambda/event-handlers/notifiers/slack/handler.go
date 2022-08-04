package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/ddb"
	"github.com/common-fate/granted-approvals/pkg/config"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/common-fate/granted-approvals/pkg/gconfig"
	slacknotifier "github.com/common-fate/granted-approvals/pkg/notifiers/slack"
	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
	"go.uber.org/zap"
)

func main() {
	var cfg config.NotificationsConfig
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

	notifier := &slacknotifier.SlackNotifier{
		DB:          db,
		FrontendURL: cfg.FrontendURL,
	}

	// @TODO
	// temporarily while making the switch to using gconfig for settings, I have implemented the slack lambda in this way.
	// In future, This should instead be implemented with some sort of registry and possibly an interface to send notifications
	// because slack currently integrates directly with the db it was hard to fit everything into the right interfaces during the initial rework
	// We will switch to a single notifications lambda which handles all configured notifications channels
	ncfg := notifier.Config()
	notificationsConfig, err := deploy.UnmarshalFeatureMap(cfg.NotificationsConfig)
	if err != nil {
		panic(err)
	}

	if slackCfg, ok := notificationsConfig[slacknotifier.NotificationsTypeSlack]; ok {
		err = ncfg.Load(ctx, &gconfig.MapLoader{Values: slackCfg})
		if err != nil {
			panic(err)
		}
		err = notifier.Init(ctx)
		if err != nil {
			panic(err)
		}
		lambda.Start(notifier.HandleEvent)
	} else {
		lambda.Start(func(ctx context.Context, event events.CloudWatchEvent) (err error) {
			log.Infow("notifications not configured, skipping handling event")
			return nil
		})
	}
}
