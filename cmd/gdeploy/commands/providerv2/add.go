package providerv2

import (
	"github.com/common-fate/common-fate/pkg/deploy"
	"github.com/common-fate/common-fate/pkg/provider"
	"github.com/common-fate/ddb"
	"github.com/urfave/cli/v2"
)

var addv2Command = cli.Command{
	Name:        "add",
	Description: "Add an access provider (development)",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "id", Usage: "An identifier for the provider"},
		&cli.StringFlag{Name: "name", Usage: "An name for the provider"},
		&cli.StringFlag{Name: "version", Usage: "An version for the provider"},
		&cli.StringFlag{Name: "url", Usage: "An url for the provider"},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		dc, err := deploy.ConfigFromContext(ctx)
		if err != nil {
			return err
		}

		//add provider to dynamo table
		outputs, err := dc.LoadOutput(ctx)
		if err != nil {
			return err
		}
		db, err := ddb.New(ctx, outputs.DynamoDBTable)
		if err != nil {
			return err
		}
		provider := provider.Provider{ID: c.String("id"), Name: c.String("name"), Version: c.String("version"), URL: c.String("url")}

		err = db.Put(ctx, &provider)
		if err != nil {
			return err
		}

		return nil
	},
}
