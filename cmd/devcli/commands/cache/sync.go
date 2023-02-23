package cache

import (
	"errors"
	"strings"

	"github.com/common-fate/common-fate/internal"
	"github.com/common-fate/common-fate/pkg/cachesync"
	"github.com/common-fate/common-fate/pkg/deploy"
	"github.com/common-fate/common-fate/pkg/pdk"
	"github.com/common-fate/common-fate/pkg/service/cachesvc"
	"github.com/common-fate/common-fate/pkg/service/requestroutersvc"
	"github.com/common-fate/ddb"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

var CacheCommand = cli.Command{
	Name:        "cache",
	Subcommands: []*cli.Command{&syncCommand},
	Action:      cli.ShowSubcommandHelp,
}

var syncCommand = cli.Command{
	Name:        "sync",
	Flags:       []cli.Flag{&cli.StringSliceFlag{Name: "deployment-mappings"}},
	Description: "Sync cache",
	Action: func(c *cli.Context) error {
		ctx := c.Context
		_ = godotenv.Load()
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
		ahc, err := internal.BuildAccessHandlerClient(ctx, internal.BuildAccessHandlerClientOpts{})
		if err != nil {
			return err
		}

		for _, dm := range c.StringSlice("deployment-mappings") {
			kv := strings.Split(dm, ":")
			if len(kv) != 2 {
				return errors.New("deployment-mapping is invalid")
			}
			pdk.LocalDeploymentMap[kv[0]] = kv[1]
		}

		// this configuration means the pdk will use the local test runtime instead of calling out to lambda
		syncer := cachesync.CacheSyncer{
			DB: db,
			Cache: cachesvc.Service{
				DB: db,
				RequestRouter: &requestroutersvc.Service{
					DB: db,
				},
				AccessHandlerClient: ahc,
			},
			AccessHandlerClient: ahc,
		}

		err = syncer.Sync(ctx)
		if err != nil {
			return err
		}

		return nil
	},
}
