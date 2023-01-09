package main

import (
	"context"
	"log"

	"github.com/common-fate/common-fate/governance/pkg/server"
	"github.com/common-fate/common-fate/pkg/config"
	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
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
	_ = godotenv.Load("../../.env")

	err := envconfig.Process(ctx, &cfg)
	if err != nil {
		return err
	}

	s, err := server.New(ctx, cfg)
	if err != nil {
		return err
	}
	return s.Start(ctx)

	// dc, err := deploy.GetDeploymentConfig()
	// if err != nil {
	// 	return err
	// }

	// log, err := logger.Build(cfg.LogLevel)
	// if err != nil {
	// 	return err
	// }
	// zap.ReplaceGlobals(log.Desugar())

	// ahc, err := internal.BuildAccessHandlerClient(ctx, internal.BuildAccessHandlerClientOpts{Region: cfg.Region, AccessHandlerURL: cfg.AccessHandlerURL})
	// if err != nil {
	// 	return err
	// }

	// api, err := api.New(ctx, api.Opts{
	// 	Log:                 log,
	// 	DynamoTable:         cfg.DynamoTable,
	// 	PaginationKMSKeyARN: cfg.PaginationKMSKeyARN,
	// 	AccessHandlerClient: ahc,
	// 	DeploymentConfig:    dc,
	// })
	// if err != nil {
	// 	return err
	// }

	// r := chi.NewRouter()
	// h := api.Handler(r)

	// host := "0.0.0.0:8889"

	// log.Infow("serving governance API", "host", host, "dynamoTable", cfg.DynamoTable)

	// return http.ListenAndServe(host, h)
}
