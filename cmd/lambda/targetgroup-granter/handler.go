package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"go.uber.org/zap"

	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/common-fate/pkg/config"
	"github.com/common-fate/common-fate/pkg/gevent"
	"github.com/common-fate/common-fate/pkg/service/requestroutersvc"
	"github.com/common-fate/common-fate/pkg/targetgroupgranter"
	"github.com/common-fate/ddb"
	"github.com/sethvargo/go-envconfig"
)

func main() {
	var cfg config.TargetGroupGranterConfig
	ctx := context.Background()
	err := envconfig.Process(ctx, &cfg)
	if err != nil {
		panic(err)
	}
	db, err := ddb.New(ctx, cfg.DynamoTable)
	if err != nil {
		panic(err)
	}
	eventBus, err := gevent.NewSender(ctx, gevent.SenderOpts{
		EventBusARN: cfg.EventBusArn,
	})
	if err != nil {
		panic(err)
	}

	granter := targetgroupgranter.Granter{
		DB: db,
		RequestRouter: &requestroutersvc.Service{
			DB: db,
		},
		EventPutter: eventBus,
	}
	log, err := logger.Build(cfg.LogLevel)
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(log.Desugar())

	lambda.Start(granter.HandleRequest)
}
