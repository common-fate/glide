package ddb

import (
	"encoding/json"
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/storage/ddbhelpers"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
)

var getUsersCommand = cli.Command{
	Name: "get-users",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "name", Aliases: []string{"n"}, Usage: "The name of the table", Required: true},
		// TODO set the region? now I user AWS_REGION=...
		// &cli.StringFlag{Name: "region", Aliases: []string{"r"}, Usage: "AWS region to provision the table into"},
		&cli.StringFlag{Name: "email", Aliases: []string{"r"}, Usage: "Search duplicates only for given email"},
		&cli.StringFlag{Name: "status", Aliases: []string{"s"}, Usage: "Status: active|deactived|all", Value: ""},
		&cli.IntFlag{Name: "limit", Aliases: []string{"l"}, Usage: "Total limit of elements to stop pagination. Set 0 for unlimited.", Value: 0},
		&cli.BoolFlag{Name: "pagination", Aliases: []string{"p"}, Usage: "Process pagination in query", Value: true},
	},

	Action: func(c *cli.Context) error {
		ctx := c.Context
		tableName := c.String("name")
		status := c.String("status")
		email := c.String("email")
		limit := c.Int("limit")
		pagination := c.Bool("pagination")

		db, err := ddb.New(ctx, tableName)
		if err != nil {
			return err
		}

		var uq ddb.QueryBuilder
		if email != "" {
			if status != "" && status != "all" {
				return fmt.Errorf("Cannot specify email=%s and status=%s at the same time", email, status)
			}
			uq = &ListUsersForEmail{Email: email}
		} else {
			switch status {
			case "", "all":
				uq = &storage.ListUsers{}
			case "active":
				uq = &storage.ListUsersForStatus{
					Status: types.IdpStatusACTIVE,
				}
			case "archived":
				uq = &storage.ListUsersForStatus{
					Status: types.IdpStatusARCHIVED,
				}
			default:
				return fmt.Errorf("Unknown status label %s", status)
			}
		}

		users := []identity.User{}
		err = ddbhelpers.QueryPages(ctx, db, uq,
			func(pageResult *ddb.QueryResult, pageQueryBuilder ddb.QueryBuilder, lastPage bool) bool {
				switch qb := pageQueryBuilder.(type) {
				case *storage.ListUsersForStatus:
					users = append(users, qb.Result...)
				case *storage.ListUsers:
					users = append(users, qb.Result...)
				case *ListUsersForEmail:
					users = append(users, qb.Result...)
				default:
					panic("Unknown type for Query Buidler")
				}
				if limit > 0 && len(users) >= limit {
					return false
				}
				return pagination
			},
		)
		if err != nil {
			return err
		}

		b, err := json.MarshalIndent(users, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(b))

		return nil
	},
}
