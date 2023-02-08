package main

import (
	"os"

	"github.com/common-fate/clio"
	"github.com/common-fate/clio/clierr"
	"github.com/common-fate/common-fate/cf/cmd/cli/commands/bootstrap"
	"github.com/common-fate/common-fate/cf/cmd/cli/commands/deployment"
	"github.com/common-fate/common-fate/cf/cmd/cli/commands/targetgroup"
	mw "github.com/common-fate/common-fate/cf/cmd/cli/middleware"
	"github.com/common-fate/common-fate/cmd/gdeploy/commands/provider"
	"github.com/common-fate/common-fate/internal/build"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

func main() {
	app := &cli.App{
		Name:        "cf",
		Description: "Common Fate CLI",
		Usage:       "Common Fate CLI",
		Version:     build.Version,
		HideVersion: false,
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "verbose", Usage: "Enable verbose logging, effectively sets environment variable GRANTED_LOG=DEBUG"},
		},
		Writer: color.Output,
		Before: func(ctx *cli.Context) error {
			if ctx.Bool("verbose") {
				os.Setenv("GRANTED_LOG", "DEBUG")
			}
			return nil
		},
		Commands: []*cli.Command{
			mw.WithBeforeFuncs(&bootstrap.Command, mw.RequireAWSCredentials()),
			mw.WithBeforeFuncs(&provider.Command, mw.RequireAWSCredentials()),
			&targetgroup.Command,
			&deployment.Command,
		},
	}

	clio.SetLevelFromEnv("CF_LOG", "GRANTED_LOG")

	zap.ReplaceGlobals(clio.G())

	err := app.Run(os.Args)
	if err != nil {
		// if the error is an instance of clio.PrintCLIErrorer then print the error accordingly
		if cliError, ok := err.(clierr.PrintCLIErrorer); ok {
			cliError.PrintCLIError()
		} else {
			clio.Error(err.Error())
		}
		os.Exit(1)
	}
}
