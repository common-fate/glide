package sync

import (
	"github.com/common-fate/granted-approvals/pkg/identity/identitysync"
	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

type SyncConfig struct {
	TableName   string `env:"APPROVALS_TABLE_NAME,default=tablename"`
	IdpProvider string `env:"IDENTITY_PROVIDER,default=COGNITO"`
	UserPoolId  string `env:"APPROVALS_COGNITO_USER_POOL_ID"`
}

var SyncCommand = cli.Command{
	Name: "sync",
	Action: func(c *cli.Context) error {
		ctx := c.Context

		var cfg SyncConfig
		_ = godotenv.Load()

		err := envconfig.Process(ctx, &cfg)
		if err != nil {
			return err
		}

		//set up the sync handler
		syncer, err := identitysync.NewIdentitySyncer(ctx, identitysync.SyncOpts{TableName: cfg.TableName, IdpType: cfg.IdpProvider, UserPoolId: cfg.UserPoolId})

		if err != nil {
			return err
		}

		zap.S().Infow("starting")
		err = syncer.Sync(ctx)
		if err != nil {
			return err
		}

		return nil
	}}
