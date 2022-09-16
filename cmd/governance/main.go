package main

import (
	"context"
	"log"
	"net/http"

	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/granted-approvals/governance"
	"github.com/common-fate/granted-approvals/internal"
	"github.com/common-fate/granted-approvals/pkg/config"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
	"go.uber.org/zap"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var cfg config.Config
	ctx := context.Background()
	_ = godotenv.Load()

	err := envconfig.Process(ctx, &cfg)
	if err != nil {
		return err
	}

	log, err := logger.Build(cfg.LogLevel)
	if err != nil {
		return err
	}
	zap.ReplaceGlobals(log.Desugar())

	ahc, err := internal.BuildAccessHandlerClient(ctx, cfg)
	if err != nil {
		return err
	}

	api, err := governance.New(ctx, governance.Opts{
		Log:                 log,
		DynamoTable:         cfg.DynamoTable,
		PaginationKMSKeyARN: cfg.PaginationKMSKeyARN,
		AccessHandlerClient: ahc,
	})
	if err != nil {
		return err
	}

	r := chi.NewRouter()
	h := api.Handler(r)

	host := "0.0.0.0:8889"

	log.Infow("serving governance API", "host", host, "dynamoTable", cfg.DynamoTable)

	return http.ListenAndServe(host, h)
}
