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
