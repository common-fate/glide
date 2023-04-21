package storage

import (
	"context"
	"os"
	"testing"

	"github.com/common-fate/ddb"
	"github.com/joho/godotenv"
)

// newTestingStorage creates a new testing storage.
// It skips the tests if the TESTING_DYNAMODB_TABLE env var isn't set.
func newTestingStorage(t *testing.T) *ddb.Client {
	ctx := context.Background()
	_ = godotenv.Load("../../.env")
	table := os.Getenv("TESTING_DYNAMODB_TABLE")
	if table == "" {
		t.Skip("TESTING_DYNAMODB_TABLE is not set")
	}
	s, err := ddb.New(ctx, table)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

// DeleteAllRequests is a helper method which can delete all requests, to be used in test cases
func deleteAllRequests(ctx context.Context, db *ddb.Client) error {
	var requests ListRequestWithGroupsWithTargets
	err := db.All(ctx, &requests)
	if err != nil {
		return err
	}
	var items []ddb.Keyer
	for r := range requests.Result {
		items = append(items, &requests.Result[r].Request)
		for g := range requests.Result[r].Groups {
			items = append(items, &requests.Result[r].Groups[g].Group)
			for t := range requests.Result[r].Groups[g].Targets {
				items = append(items, &requests.Result[r].Groups[g].Targets[t])
			}
		}
	}

	return db.DeleteBatch(ctx, items...)
}
