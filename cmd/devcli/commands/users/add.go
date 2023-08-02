package users

import (
	"fmt"

	"github.com/benbjohnson/clock"
	"github.com/common-fate/common-fate/pkg/config"
	"github.com/common-fate/common-fate/pkg/deploy"
	"github.com/common-fate/common-fate/pkg/gconfig"
	"github.com/common-fate/common-fate/pkg/identity/identitysync"
	"github.com/common-fate/common-fate/pkg/service/cognitosvc"
	"github.com/common-fate/ddb"
	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
	"github.com/urfave/cli/v2"
)

// create a lot of users in cognito
var UsersCommand = cli.Command{
	Name:        "users",
	Subcommands: []*cli.Command{&addCommand},
	Action:      cli.ShowSubcommandHelp,
}

var addCommand = cli.Command{
	Name:        "add",
	Flags:       []cli.Flag{},
	Description: "Add a Cognito user to a group",
	Action: func(c *cli.Context) error {
		ctx := c.Context
		var cfg config.Config
		_ = godotenv.Load()

		err := envconfig.Process(ctx, &cfg)
		if err != nil {
			return err
		}

		tokenizer, err := ddb.NewKMSTokenizer(ctx, cfg.PaginationKMSKeyARN)
		if err != nil {
			return err
		}

		db, err := ddb.New(ctx, cfg.DynamoTable, ddb.WithPageTokenizer(tokenizer))
		if err != nil {
			return err
		}

		clk := clock.New()

		if err != nil {
			return err
		}

		cog := &identitysync.CognitoSync{}
		err = cog.Config().Load(ctx, &gconfig.MapLoader{Values: map[string]string{"userPoolId": cfg.CognitoUserPoolID}})
		if err != nil {
			return err
		}
		err = cog.Init(ctx)
		if err != nil {
			return err
		}
		ic, err := deploy.UnmarshalFeatureMap(cfg.IdentitySettings)
		if err != nil {
			return err
		}

		idsync, err := identitysync.NewIdentitySyncer(ctx, identitysync.SyncOpts{
			TableName:           cfg.DynamoTable,
			UserPoolId:          cfg.CognitoUserPoolID,
			IdpType:             cfg.IdpProvider,
			IdentityConfig:      ic,
			IdentityGroupFilter: cfg.IdentityGroupFilter,
		})
		if err != nil {
			return err
		}

		Cognito := &cognitosvc.Service{
			Clock:        clk,
			DB:           db,
			Syncer:       idsync,
			Cognito:      cog,
			AdminGroupID: cfg.AdminGroup,
		}

		for i := 50; i <= 100; i++ {

			createUser := cognitosvc.CreateUserOpts{

				FirstName: fmt.Sprintf("User%d", i),
				LastName:  "Doe",
				Email:     fmt.Sprintf("user%d@commonfate.io", i),
				IsAdmin:   false,
			}
			_, err := Cognito.AdminCreateUser(ctx, createUser)
			if err != nil {
				return err
			}
		}

		return nil
	},
}
