package main

import (
	"context"
	"encoding/json"

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
	ncfg := notifier.Config()
	var notificationsConfig []deploy.Feature
	err = json.Unmarshal([]byte(cfg.NotificationsConfig), &notificationsConfig)
	if err != nil {
		panic(err)
	}
	var found bool
	for _, c := range notificationsConfig {
		if c.Uses == "commonfate/notifications/slack@v1" {
			found = true
			b, err := json.Marshal(c.With)
			if err != nil {
				panic(err)
			}
			err = ncfg.Load(ctx, gconfig.JSONLoader{Data: b})
			if err != nil {
				panic(err)
			}
			err = notifier.Init(ctx)
			if err != nil {
				panic(err)
			}
		}
	}
	if !found {
		lambda.Start(func(ctx context.Context, event events.CloudWatchEvent) (err error) {
			log.Infow("notifications not configured, skipping handling event")
			return nil
		})
	} else {
		lambda.Start(notifier.HandleEvent)
	}

}
