package providerv2

import (
	"github.com/common-fate/common-fate/pkg/deploy"
	"github.com/common-fate/common-fate/pkg/provider"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
	"github.com/urfave/cli/v2"
)

var addv2Command = cli.Command{
	Name:        "add",
	Description: "Add an access provider (development)",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "alias", Usage: "An alias for the provider"},
		&cli.StringFlag{Name: "function-arn", Usage: "An arn for the provider lambda function"},
		&cli.StringFlag{Name: "icon-name", Usage: "A type for the provider"},
		&cli.StringFlag{Name: "team", Usage: "the vendor team"},
		&cli.StringFlag{Name: "name", Usage: "the registry name"},
		&cli.StringFlag{Name: "version", Usage: "the version for the provider"},
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
		funcArn := c.String("function-arn")
		provider := provider.Provider{ID: types.NewProviderID(), Team: c.String("team"), Name: c.String("name"), Version: c.String("version"), IconName: c.String("icon-name"), FunctionARN: &funcArn, Alias: c.String("alias")}

		err = db.Put(ctx, &provider)
		if err != nil {
			return err
		}

		return nil
	},
}
