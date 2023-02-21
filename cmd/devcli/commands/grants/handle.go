package grants

import (
	"time"

	"github.com/common-fate/clio"
	"github.com/common-fate/common-fate/accesshandler/pkg/types"
	ahTypes "github.com/common-fate/common-fate/accesshandler/pkg/types"
	"github.com/common-fate/common-fate/pkg/config"
	"github.com/common-fate/common-fate/pkg/service/requestroutersvc"
	"github.com/common-fate/common-fate/pkg/targetgroupgranter"
	"github.com/common-fate/ddb"
	"github.com/common-fate/iso8601"
	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
	"github.com/urfave/cli/v2"
)

var Handle = cli.Command{
	Name: "handle",
	Action: func(c *cli.Context) error {
		ctx := c.Context
		var cfg config.TargetGroupGranterConfig
		_ = godotenv.Load("../../../../.env")
		err := envconfig.Process(ctx, &cfg)
		if err != nil {
			panic(err)
		}
		db, err := ddb.New(ctx, cfg.DynamoTable)
		if err != nil {
			panic(err)
		}
		granter := targetgroupgranter.Granter{
			Cfg: cfg,
			DB:  db,
			RequestRouter: &requestroutersvc.Service{
				DB: db,
			},
		}

		out, err := granter.HandleRequest(ctx, targetgroupgranter.InputEvent{
			Grant: types.Grant{
				Subject:  openapi_types.Email("josh@commonfate.io"),
				Start:    iso8601.New(time.Now().Add(time.Second * 2)),
				End:      iso8601.New(time.Now().Add(time.Hour)),
				Provider: "josh-example",
				ID:       ahTypes.NewGrantID(),
				With: ahTypes.Grant_With{
					AdditionalProperties: map[string]string{
						"vault": "test",
					},
				},
			},
		})
		if err != nil {
			return err
		}
		clio.Infow("complete", "out", out)
		return nil
	},
}
