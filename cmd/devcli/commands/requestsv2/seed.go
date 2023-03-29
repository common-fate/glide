package requestsv2

import (
	"github.com/common-fate/common-fate/pkg/config"
	"github.com/common-fate/common-fate/pkg/requestsv2.go"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
	"github.com/urfave/cli/v2"
)

var SeedCommand = cli.Command{
	Name:        "seed",
	Description: "Seeds some dummy data into the dynamo table for testing new workflows",
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

		items := []ddb.Keyer{}

		//create an entitlement
		ent := requestsv2.Entitlement{
			ID: types.NewEntitlementID(),
			Kind: requestsv2.TargetFrom{
				Kind:      "Account",
				Name:      "AWS",
				Publisher: "common-fate",
				Version:   "v0.1.0",
			},
			OptionSchema: types.TargetSchema{
				AdditionalProperties: map[string]types.TargetArgument{
					"accountId":        {Title: "Account"},
					"permissionSetArn": {Title: "Permission Set"},
				},
			},
		}

		ent2 := requestsv2.Entitlement{
			ID: types.NewEntitlementID(),

			Kind: requestsv2.TargetFrom{
				Kind:      "Group",
				Name:      "Okta",
				Publisher: "common-fate",
				Version:   "v0.1.0",
			},
			OptionSchema: types.TargetSchema{
				AdditionalProperties: map[string]types.TargetArgument{
					"groupName": {Title: "Group Name"},
				},
			},
		}

		items = append(items, &ent)
		items = append(items, &ent2)

		//create some options
		opt1 := requestsv2.Option{

			Label: "accountId",
			Value: "123456789012",
			Provider: requestsv2.TargetFrom{
				Kind:      "Account",
				Name:      "AWS",
				Publisher: "common-fate",
				Version:   "v0.1.0",
			},
		}
		opt2 := requestsv2.Option{
			Label: "permissionSet",
			Value: "123-abc",
			Provider: requestsv2.TargetFrom{
				Kind:      "Account",
				Name:      "AWS",
				Publisher: "common-fate",
				Version:   "v0.1.0",
			},
		}

		opt3 := requestsv2.Option{
			Label: "groupName",
			Value: "This is a okta group",
			Provider: requestsv2.TargetFrom{
				Kind:      "Group",
				Name:      "Okta",
				Publisher: "common-fate",
				Version:   "v0.1.0",
			},
		}
		items = append(items, &opt1)
		items = append(items, &opt2)
		items = append(items, &opt3)

		err = db.PutBatch(ctx, items...)
		if err != nil {
			return err
		}
		return nil

	}),
}
