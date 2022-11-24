package main

import (
	"fmt"
	"os"

	"github.com/common-fate/granted-approvals/accesshandler/cmd/cli/commands/grants"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

func main() {
	flags := []cli.Flag{
		&cli.StringFlag{Name: "api-url", Value: "http://localhost:9092", EnvVars: []string{"ACCESS_HANDLER_URL"}, Hidden: true},
	}

	app := &cli.App{
		Flags:                flags,
		Name:                 "cli",
		UsageText:            "cli [global options] command [command options] [arguments...]",
		HideVersion:          false,
		Commands:             []*cli.Command{&grants.Command},
		EnableBashCompletion: true,
	}

	logCfg := zap.NewDevelopmentConfig()
	logCfg.DisableStacktrace = true

	log, err := logCfg.Build()
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
	zap.ReplaceGlobals(log)

	err = app.Run(os.Args)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
}
