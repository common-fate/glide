package commands

import (
	"github.com/common-fate/common-fate/internal"
	"github.com/common-fate/common-fate/pkg/cachesync"
	"github.com/common-fate/common-fate/pkg/config"
	"github.com/common-fate/common-fate/pkg/deploy"
	"github.com/common-fate/common-fate/pkg/service/cachesvc"
	"github.com/common-fate/ddb"
	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
	"github.com/urfave/cli/v2"
)

var CacheCommand = cli.Command{
	Name:        "cache",
	Subcommands: []*cli.Command{&syncCommand},
	Action:      cli.ShowSubcommandHelp,
}

var syncCommand = cli.Command{
	Name:        "sync",
	Flags:       []cli.Flag{},
	Description: "Sync schemas from PDK",
	Action: func(c *cli.Context) error {
		ctx := c.Context
		var cfg config.CacheSyncConfig
		_ = godotenv.Load()

		err := envconfig.Process(ctx, &cfg)
		if err != nil {
			return err
		}
		do, err := deploy.LoadConfig(deploy.DefaultFilename)
		if err != nil {
			return err
		}
		o, err := do.LoadOutput(ctx)
		if err != nil {
			return err
		}

		db, err := ddb.New(ctx, o.DynamoDBTable)
		if err != nil {
			return err
		}

		ahc, err := internal.BuildAccessHandlerClient(ctx, internal.BuildAccessHandlerClientOpts{Region: cfg.Region, AccessHandlerURL: cfg.AccessHandlerURL})
		if err != nil {
			panic(err)
		}
		dc, err := deploy.GetDeploymentConfig()
		if err != nil {
			return err
		}
		syncer := cachesync.CacheSyncer{
			DB:                  db,
			AccessHandlerClient: ahc,
			Cache: cachesvc.Service{
				DB:                   db,
				AccessHandlerClient:  ahc,
				ProviderConfigReader: dc,
			},
		}

		err = syncer.Sync(ctx)
		if err != nil {
			return err
		}

		return nil
	},
}
