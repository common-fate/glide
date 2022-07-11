package sync

import (
	"github.com/common-fate/granted-approvals/pkg/clio"
	"github.com/common-fate/granted-approvals/pkg/config"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/common-fate/granted-approvals/pkg/identity/identitysync"
	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
	"github.com/urfave/cli/v2"
)

var SyncCommand = cli.Command{
	Name: "sync",
	Action: func(c *cli.Context) error {
		ctx := c.Context

		var cfg config.SyncConfig
		_ = godotenv.Load()

		err := envconfig.Process(ctx, &cfg)
		if err != nil {
			return err
		}

		ic, err := deploy.UnmarshalIdentity(cfg.IdentitySettings)
		if err != nil {
			panic(err)
		}

		//set up the sync handler
		syncer, err := identitysync.NewIdentitySyncer(ctx, identitysync.SyncOpts{
			TableName:      cfg.TableName,
			IdpType:        cfg.IdpProvider,
			UserPoolId:     cfg.UserPoolId,
			IdentityConfig: ic,
		})

		if err != nil {
			return err
		}

		clio.Info("Starting sync")
		err = syncer.Sync(ctx)
		if err != nil {
			return err
		}

		return nil
	}}
