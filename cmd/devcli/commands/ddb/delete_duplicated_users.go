package ddb

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/storage/ddbhelpers"
	"github.com/common-fate/ddb"
)

type DuplicateUsersReport struct {
	DuplicatesToDelete   map[string][]identity.User
	InitialToKeep        map[string]identity.User
	TotalDuplicatesCount int
}

func findDuplicates(users []identity.User, maxDups int) *DuplicateUsersReport {
	result := DuplicateUsersReport{
		DuplicatesToDelete: map[string][]identity.User{},
		InitialToKeep:      map[string]identity.User{},
	}

	for _, u := range users {
		if u2, ok := result.InitialToKeep[u.Email]; !ok {
			result.InitialToKeep[u.Email] = u
		} else {
			if _, ok := result.DuplicatesToDelete[u.Email]; !ok {
				result.DuplicatesToDelete[u.Email] = []identity.User{}
			}
			if u2.CreatedAt.Before(u.CreatedAt) {
				result.DuplicatesToDelete[u.Email] = append(result.DuplicatesToDelete[u.Email], u)
			} else {
				result.InitialToKeep[u.Email] = u2
				result.DuplicatesToDelete[u.Email] = append(result.DuplicatesToDelete[u.Email], u2)
			}
			result.TotalDuplicatesCount = result.TotalDuplicatesCount + 1
			if maxDups > 0 && result.TotalDuplicatesCount >= maxDups {
				break
			}
		}
	}
	return &result
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
		&cli.IntFlag{Name: "workers", Usage: "Concurrent workers deleting entries.", Value: 1},
	},

	Action: func(c *cli.Context) error {
		ctx, _ := context.WithCancel(c.Context)
		tableName := c.String("name")
		email := c.String("email")
		limit := c.Int("limit")
		pagination := c.Bool("pagination")
		dryRun := c.Bool("dry-run")
		maxDups := c.Int("max-dups")
		numWorkers := c.Int("workers")

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

		fmt.Fprintf(os.Stderr, "Found %d duplicated users\n", duplicatedUsers.TotalDuplicatesCount)
		for k, v := range duplicatedUsers.DuplicatesToDelete {
			fmt.Fprintf(os.Stderr, " - %s: %d\n", k, len(v))
		}

		if !dryRun {
			start := time.Now()
			fmt.Fprintf(os.Stderr, "Deleting duplicates at %s using %d workers...\n", start.Format(time.RFC3339), numWorkers)

			ch := make(chan []identity.User)
			var wg sync.WaitGroup

			for i := 0; i < numWorkers; i++ {
				go func() {
					defer wg.Done()
					wg.Add(1)
					for dups := range ch {
						email := dups[0].Email
						start := time.Now()
						fmt.Fprintf(os.Stderr, "Deleting %d dups for email=%s at %s...\n", len(dups), email, start.Format(time.RFC3339))
						entriesToDelete := make([]ddb.Keyer, len(dups))
						for i, u := range dups {
							entriesToDelete[i] = &UserBaseKeyOnly{ID: u.ID}
						}
						err = db.DeleteBatch(ctx, entriesToDelete...)
						if err != nil {
							fmt.Fprintf(os.Stderr, "Error processing %s: %s\n", email, err)
							break
						}
						end := time.Now()
						elapsed := end.Sub(start)
						fmt.Fprintf(os.Stderr, "Done deleting dups for email=%s Duration=%s...\n", email, elapsed)
					}
				}()
			}
			for _, dups := range duplicatedUsers.DuplicatesToDelete {
				ch <- dups
			}
			close(ch)
			wg.Wait()

			end := time.Now()
			elapsed := end.Sub(start)
			fmt.Fprintf(os.Stderr, "Deletion complete at %s. Duration %s...\n", end.Format(time.RFC3339), elapsed)
		} else {
			fmt.Fprintln(os.Stderr, "WARNING: Dry-run mode, skipping deletion...")
		}

		return nil
	},
}
