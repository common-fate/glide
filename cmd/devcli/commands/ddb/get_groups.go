package ddb

import (
	"encoding/json"
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/storage/ddbhelpers"
	"github.com/common-fate/ddb"
)

var getGroupsCommand = cli.Command{
	Name: "get-groups",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "name", Aliases: []string{"n"}, Usage: "The name of the table", Required: true},
		// TODO set the region? now I user AWS_REGION=...
		// &cli.StringFlag{Name: "region", Aliases: []string{"r"}, Usage: "AWS region to provision the table into"},
		&cli.IntFlag{Name: "limit", Aliases: []string{"l"}, Usage: "Total limit of elements to stop pagination. Set 0 for unlimited.", Value: 0},
		&cli.BoolFlag{Name: "pagination", Aliases: []string{"p"}, Usage: "Process pagination in query", Value: true},
	},

	Action: func(c *cli.Context) error {
		ctx := c.Context
		tableName := c.String("name")
		limit := c.Int("limit")
		pagination := c.Bool("pagination")

		db, err := ddb.New(ctx, tableName)
		if err != nil {
			return err
		}

		gq := &storage.ListGroups{}

		var groups []identity.Group
		err = ddbhelpers.QueryPages(ctx, db, gq,
			func(pageResult *ddb.QueryResult, pageQueryBuilder ddb.QueryBuilder, lastPage bool) bool {
				if qb, ok := pageQueryBuilder.(*storage.ListGroups); ok {
					groups = append(groups, qb.Result...)
				} else {
					panic("Unknown type for QueryBuilder")
				}
				if limit > 0 && len(groups) >= limit {
					return false
				}
				return pagination
			},
		)
		if err != nil {
			return err
		}

		b, err := json.Marshal(gq.Result)
		if err != nil {
			return err
		}
		fmt.Println(string(b))

		return nil
	},
}
