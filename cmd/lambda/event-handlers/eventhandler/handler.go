package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/benbjohnson/clock"

	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/common-fate/pkg/config"
	"github.com/common-fate/common-fate/pkg/eventhandler"
	"github.com/common-fate/common-fate/pkg/gevent"
	"github.com/common-fate/common-fate/pkg/service/requestroutersvc"
	"github.com/common-fate/common-fate/pkg/service/workflowsvc"
	"github.com/common-fate/common-fate/pkg/service/workflowsvc/runtimes/live"

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
	eb, err := gevent.NewSender(ctx, gevent.SenderOpts{
		EventBusARN: cfg.EventBusArn,
	})

	if err != nil {
		panic(err)
	}

	clk := clock.New()
	eventHandler := eventhandler.EventHandler{
		DB:       db,
		Eventbus: eb,
		Workflow: &workflowsvc.Service{
			DB:       db,
			Clk:      clk,
			Eventbus: eb,
			Runtime: &live.Runtime{
				StateMachineARN: cfg.StateMachineARN,
				Eventbus:        eb,
				DB:              db,
				RequestRouter: &requestroutersvc.Service{
					DB: db,
				},
			},
		},
	}
	log, err := logger.Build(cfg.LogLevel)
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(log.Desugar())
	zap.S().Infow("starting event handler with configuration", "config", cfg)
	lambda.Start(eventHandler.HandleEvent)
}
