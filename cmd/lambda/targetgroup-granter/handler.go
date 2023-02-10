package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/common-fate/common-fate/pkg/config"
	"github.com/common-fate/common-fate/pkg/service/requestroutersvc"
	"github.com/common-fate/common-fate/pkg/targetgroupgranter"
	"github.com/common-fate/ddb"
	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

func main() {
	var cfg config.TargetGroupGranterConfig
	ctx := context.Background()
	_ = godotenv.Load()
	ctx.Deadline()

	err := envconfig.Process(ctx, &cfg)
	if err != nil {
		panic(err)
	}
	db, err := ddb.New(ctx, cfg.DynamoTable)
	if err != nil {
		panic(err)
	}
	granter := targetgroupgranter.Granter{
		Cfg: cfg,
		DB:  db,
		RequestRouter: &requestroutersvc.Service{
			DB: db,
		},
	}

	lambda.Start(granter.HandleRequest)
}
