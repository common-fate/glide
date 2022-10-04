package storage

import (
	"testing"

	"github.com/common-fate/ddb"
	"github.com/common-fate/ddb/ddbtest"
	"github.com/common-fate/granted-approvals/pkg/cache"
)

func TestGetProviderOptions(t *testing.T) {
	db := newTestingStorage(t)

	po := cache.ProviderOption{
		Provider: "test",
		Arg:      "test",
		Label:    "test",
		Value:    "test",
	}
	ddbtest.PutFixtures(t, db, &po)

	tc := []ddbtest.QueryTestCase{
		{
			Name:  "ok",
			Query: &ListProviderOptionsForArg{ProviderID: "test", ArgID: "test"},
			Want:  &ListProviderOptionsForArg{ProviderID: "test", ArgID: "test", Result: []cache.ProviderOption{po}},
		},
		{
			Name:    "not found",
			Query:   &ListProviderOptionsForArg{ProviderID: "somethingelse", ArgID: "test"},
			WantErr: ddb.ErrNoItems,
		},
	}

	ddbtest.RunQueryTests(t, db, tc)
}
