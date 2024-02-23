package main

import (
	"fmt"
	"os"

	"github.com/common-fate/common-fate/cmd/devcli/commands/cache"
	"github.com/common-fate/common-fate/cmd/devcli/commands/db"
	"github.com/common-fate/common-fate/cmd/devcli/commands/ddb"
	"github.com/common-fate/common-fate/cmd/devcli/commands/events"
	"github.com/common-fate/common-fate/cmd/devcli/commands/grants"
	"github.com/common-fate/common-fate/cmd/devcli/commands/groups"
	registry "github.com/common-fate/common-fate/cmd/devcli/commands/provider-registry"
	"github.com/common-fate/common-fate/cmd/devcli/commands/slack"

	"github.com/common-fate/common-fate/cmd/devcli/commands/healthcheck"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

func main() {
	app := &cli.App{
		Name:        "commonfate",
		Writer:      color.Output,
		Version:     "v0.0.1",
		HideVersion: false,
		Commands: []*cli.Command{
			&groups.GroupsCommand,
			&db.DBCommand,
			&ddb.DDBCommand,
			&events.EventsCommand,
			&slack.SlackCommand,
			&cache.CacheCommand,
			&healthcheck.Command,
			&grants.Command,
			&registry.Command,
		},
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
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
