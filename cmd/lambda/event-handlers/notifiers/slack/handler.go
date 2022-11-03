package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
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

	h := handler{
		Log:         log,
		DB:          db,
		FrontendURL: cfg.FrontendURL,
	}

	lambda.Start(h.handleEvent)
}

type handler struct {
	Log         *zap.SugaredLogger
	DB          ddb.Storage
	FrontendURL string
}

func (h *handler) handleEvent(ctx context.Context, event events.CloudWatchEvent) error {
	notifier := &slacknotifier.SlackNotifier{
		DB:          h.DB,
		FrontendURL: h.FrontendURL,
	}

	dc, err := deploy.GetDeploymentConfig()
	if err != nil {
		return err
	}

	// don't cache notification config - re-read it every time the Lambda executes.
	// This avoids us using stale config if we're reading config from a remote API,
	// rather than from env vars. This adds latency but this is an async operation
	// anyway so it doesn't really matter.
	notificationsConfig, err := dc.ReadNotifications(ctx)
	if err != nil {
		h.Log.Errorw("failed to initialise slack notifier", "error", err)
		return err
	}

	err = notifier.Init(ctx, notificationsConfig)
	if err != nil {
		h.Log.Errorw("failed to initialise slack notifier", "error", err)
		return err
	}
	return notifier.HandleEvent(ctx, event)
}
