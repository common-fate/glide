package targetgroup

import (
	"errors"

	"github.com/common-fate/common-fate/pkg/config"
	"github.com/common-fate/common-fate/pkg/target"
	"github.com/common-fate/ddb"
	"github.com/common-fate/provider-registry-sdk-go/pkg/handlerclient"
	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
	"github.com/urfave/cli/v2"
)

// this command can be run in dev with:
// go run cf/cmd/cli/main.go healthcheck

var Command = cli.Command{
	Name:        "targetgroup",
	Description: "manage a target group",
	Usage:       "manage a target group",
	Subcommands: []*cli.Command{
		&CreateCommand,
	},
}

var CreateCommand = cli.Command{
	Name:        "create",
	Description: "create a target group using a local provider repo and without connecting to any provider registry",
	Usage:       "create a target group using a local provider repo and without connecting to any provider registry",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "path", Required: true},
		&cli.StringFlag{Name: "id", Required: true},
		&cli.StringFlag{Name: "kind", Usage: "the target kind that the provider grants access to", Required: true},
		&cli.StringFlag{Name: "provider", Usage: "publisher/name@version", Required: true},
	},

	Action: cli.ActionFunc(func(c *cli.Context) error {
		ctx := c.Context
		// Read from the .env file
		var cfg config.HealthCheckerConfig
		_ = godotenv.Load()
		err := envconfig.Process(ctx, &cfg)
		if err != nil {
			return err
		}
		db, err := ddb.New(ctx, cfg.TableName)
		if err != nil {
			return err
		}

		hc := handlerclient.Client{Executor: handlerclient.Local{Dir: c.String("path")}}

		describe, err := hc.Describe(ctx)
		if err != nil {
			return err
		}
		if describe.Schema.Targets == nil {
			return errors.New("target schema was nil")
		}
		tg := target.Group{
			ID:   c.String("id"),
			From: target.From{Publisher: c.String("provider")},
			// FIXME:
			// Schema: (*describe.Schema.Targets)[c.String("kind")],
		}
		err = db.Put(ctx, &tg)
		if err != nil {
			return err
		}
		return nil
	}),
}
