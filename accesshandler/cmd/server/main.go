package main

import (
	"context"
	"log"

	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/config"
	"go.uber.org/zap"

	"github.com/common-fate/granted-approvals/accesshandler/pkg/server"
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
	log.Warn("setting AssumeExecutionRoleARN to '' because the access handler is running in development mode. This means that AWS api calls will be made using your current credentials without assuming the access handler role. It is safe to ignore this warning.")
	cfg.AssumeExecutionRoleARN = ""

	s, err := server.New(ctx, cfg)
	if err != nil {
		return err
	}

	return s.Start(ctx)
}
