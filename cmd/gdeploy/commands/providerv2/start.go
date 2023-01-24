package providerv2

import (
	"context"

	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/common-fate/pkg/config"
	"github.com/common-fate/common-fate/pkg/deploymentcli/server"

	"github.com/joho/godotenv"
	"github.com/pkg/browser"
	"github.com/sethvargo/go-envconfig"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var StartCommand = cli.Command{
	Name:        "start",
	Description: "Start a new deployment management session",
	Usage:       "Start a new deployment management session",
	Flags:       []cli.Flag{},
	Action: func(c *cli.Context) error {
		return run()
	},
}

func run() error {
	var cfg config.ProviderDeploymentCLI
	ctx := context.Background()
	err := godotenv.Load()
	if err != nil {
		return err
	}
	err = envconfig.Process(ctx, &cfg)
	if err != nil {
		return err
	}

	log, err := logger.Build(cfg.LogLevel)
	if err != nil {
		return err
	}
	zap.ReplaceGlobals(log.Desugar())

	s, err := server.New(ctx, server.Opts{
		Logger: log,
		Cfg:    cfg,
	})
	if err != nil {
		return err
	}

	log.Infow("starting server", "config", cfg)
	go browser.OpenURL(cfg.LocalFrontendURL)
	return s.Start(ctx)
}
