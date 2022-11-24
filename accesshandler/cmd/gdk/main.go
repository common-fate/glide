package main

import (
	"fmt"
	"os"

	"github.com/common-fate/granted-approvals/accesshandler/cmd/gdk/commands/fixtures"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

func main() {
	app := &cli.App{
		Name:                 "gdk",
		Usage:                "https://granted.dev",
		Description:          "Common Fate Development Kit",
		UsageText:            "gdk [global options] command [command options] [arguments...]",
		HideVersion:          false,
		Commands:             []*cli.Command{&fixtures.Command},
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
