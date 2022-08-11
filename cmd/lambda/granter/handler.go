package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/config"
	lambdagranter "github.com/common-fate/granted-approvals/accesshandler/pkg/runtime/lambda/granter"
	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

func main() {
	var cfg config.GranterConfig
	ctx := context.Background()
	_ = godotenv.Load()
	ctx.Deadline()

	err := envconfig.Process(ctx, &cfg)
	if err != nil {
		panic(err)
	}

	g, err := lambdagranter.NewGranter(ctx, cfg)
	if err != nil {
		panic(err)
	}

	lambda.Start(g.HandleRequest)
}
