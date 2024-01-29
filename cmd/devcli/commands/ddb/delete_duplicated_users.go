package ddb

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/urfave/cli/v2"

	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/storage/ddbhelpers"
	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/ddb"
)

// Needed to filter by duplicates. Not needed in normal operations
type ListUsersForEmail struct {
	Result []identity.User `ddb:"result"`
	Email  string
}

func (l *ListUsersForEmail) BuildQuery() (*dynamodb.QueryInput, error) {

	qi := dynamodb.QueryInput{
		IndexName:              aws.String(keys.IndexNames.GSI2),
		KeyConditionExpression: aws.String("GSI2PK = :pk2 and GSI2SK = :sk2"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk2": &types.AttributeValueMemberS{Value: keys.Users.GSI2PK},
			":sk2": &types.AttributeValueMemberS{Value: keys.Users.GSI2SK(string(l.Email))},
		},
	}
	return &qi, nil
}

func findDuplicates(users []identity.User, maxDups int) []identity.User {
	toDelete := []identity.User{}
	userMap := map[string]identity.User{}

	for _, u := range users {
		if u2, ok := userMap[u.Email]; ok {
			if u2.CreatedAt.Before(u.CreatedAt) {
				toDelete = append(toDelete, u)
			} else {
				userMap[u.Email] = u
				toDelete = append(toDelete, u2)
			}
		} else {
			userMap[u.Email] = u
		}
		if maxDups > 0 && len(toDelete) >= maxDups {
			break
		}
	}
	return toDelete
}

var deleteDuplicatedUsersCommand = cli.Command{
	Name: "delete-duplicated-users",
	Description: `Find any duplicated users in the Dynamodb DB.

This command fixes the issue of duplicated users added during sync due
the lack of pagination querying dynamodb.

It will list all users in dynamodb, identifying those with a duplicated
email, and keeping the one with the newest creation date.

Runs by default in dry-run mode, or allows to delete them from dynamodb.
	`,
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "name", Usage: "The name of the table", Required: true},
		// TODO set the region? now I user AWS_REGION=...
		// &cli.StringFlag{Name: "region", Aliases: []string{"r"}, Usage: "AWS region to provision the table into"},
		&cli.StringFlag{Name: "email", Aliases: []string{"r"}, Usage: "Search duplicates only for given email"},
		&cli.IntFlag{Name: "limit", Usage: "Total limit of elements to stop pagination. Set 0 for unlimited.", Value: 0},
		&cli.BoolFlag{Name: "pagination", Usage: "Process pagination in query", Value: true},
		&cli.BoolFlag{Name: "dry-run", Usage: "Set to false to actually delete the users. Default true", Value: true},
		&cli.IntFlag{Name: "max-dups", Usage: "Only delete a max of given duplicates", Value: 0},
	},

	Action: func(c *cli.Context) error {
		ctx := c.Context
		tableName := c.String("name")
		email := c.String("email")
		limit := c.Int("limit")
		pagination := c.Bool("pagination")
		dryRun := c.Bool("dry-run")
		maxDups := c.Int("max-dups")

		db, err := ddb.New(ctx, tableName)
		if err != nil {
			return err
		}

		var uq ddb.QueryBuilder
		if email == "" {
			uq = &storage.ListUsers{}
		} else {
			uq = &ListUsersForEmail{Email: email}
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

		fmt.Fprintf(os.Stderr, "Got %d users\n", len(users))
		fmt.Fprintln(os.Stderr, "Calculating duplicates...")

		duplicatedUsers := findDuplicates(users, maxDups)

		b, err := json.MarshalIndent(duplicatedUsers, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(b))

		fmt.Fprintf(os.Stderr, "Found %d duplicated\n", len(duplicatedUsers))

		if !dryRun {
			fmt.Fprintln(os.Stderr, "Deleting duplicates...")
			entriesToDelete := make([]ddb.Keyer, len(duplicatedUsers))
			for i, u := range duplicatedUsers {
				entriesToDelete[i] = &u
			}
			err = db.DeleteBatch(ctx, entriesToDelete...)
			if err != nil {
				return err
			}
		} else {
			fmt.Fprintln(os.Stderr, "WARNING: Dry-run mode, skipping deletion...")
		}

		return nil
	},
}
