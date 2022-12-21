package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/common-fate/governance"
	"github.com/common-fate/common-fate/internal"
	"github.com/common-fate/common-fate/pkg/config"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
	"go.uber.org/zap"
)

func buildHandler() (http.Handler, error) {
	var cfg config.Config
	ctx := context.Background()
	_ = godotenv.Load("../../.env")

	err := envconfig.Process(ctx, &cfg)
	if err != nil {
		panic(err)
	}

	log, err := logger.Build(cfg.LogLevel)
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(log.Desugar())

	ahc, err := internal.BuildAccessHandlerClient(ctx, internal.BuildAccessHandlerClientOpts{Region: cfg.Region, AccessHandlerURL: cfg.AccessHandlerURL})
	if err != nil {
		panic(err)
	}

	api, err := governance.New(ctx, governance.Opts{
		Log:                 log,
		DynamoTable:         cfg.DynamoTable,
		PaginationKMSKeyARN: cfg.PaginationKMSKeyARN,
		AccessHandlerClient: ahc,
	})
	if err != nil {
		panic(err)
	}
	r := chi.NewRouter()

	return api.Handler(r), nil
}

func main() {
	l, err := buildHandler()
	if err != nil {
		panic(err)
	}

	lambda.Start(l)
}
