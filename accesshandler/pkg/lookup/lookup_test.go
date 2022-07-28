package lookup

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers/azure/ad"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers/okta"
	"github.com/stretchr/testify/assert"
)

var testRegistry = ProviderRegistry{
	Providers: map[string]RegisteredProvider{
		"commonfate/okta@v1": {
			Provider:    &okta.Provider{},
			DefaultID:   "okta",
			Description: "Okta groups",
		},
		"commonfate/azure-ad@v1": {
			Provider:    &ad.Provider{},
			DefaultID:   "azure-ad",
			Description: "Azure-AD groups",
		},
	},
}

func TestFromCLIOption(t *testing.T) {
	type testcase struct {
		name    string
		give    string
		want    RegisteredProvider
		wantKey string
		wantErr error
	}

	testcases := []testcase{
		{
			name:    "ok okta",
			give:    "Okta groups (commonfate/okta@v1)",
			wantKey: "commonfate/okta@v1",
			want:    testRegistry.Providers["commonfate/okta@v1"],
		},
		{
			name:    "ok azure",
			give:    "Azure-AD groups (commonfate/azure-ad@v1)",
			wantKey: "commonfate/azure-ad@v1",
			want:    testRegistry.Providers["commonfate/azure-ad@v1"],
		},
		{
			name:    "from CLIOptions okta",
			give:    testRegistry.CLIOptions()[0],
			wantKey: "commonfate/okta@v1",
			want:    testRegistry.Providers["commonfate/okta@v1"],
		},
		{
			name:    "from CLIOptions azure",
			give:    testRegistry.CLIOptions()[1],
			wantKey: "commonfate/azure-ad@v1",
			want:    testRegistry.Providers["commonfate/azure-ad@v1"],
		},
		{
			name:    "invalid format opkta",
			give:    "commonfate/okta@v1",
			wantErr: errors.New("couldn't extract provider key: commonfate/okta@v1"),
		},
		{
			name:    "invalid format azure",
			give:    "commonfate/azure-ad@v1",
			wantErr: errors.New("couldn't extract provider key: commonfate/azure-ad@v1"),
		},
		{
			name:    "provider not found",
			give:    "Test Provider (commonfate/something-else@v1)",
			wantErr: errors.New("couldn't find provider with key: commonfate/something-else@v1"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			gotKey, got, err := testRegistry.FromCLIOption(tc.give)
			if err != nil && tc.wantErr == nil {
				t.Fatal(err)
			}
			if tc.wantErr != nil {
				assert.Equal(t, tc.wantErr, err)
			}
			assert.Equal(t, tc.want, got)
			assert.Equal(t, tc.wantKey, gotKey)
		})
	}
}

// The following test enforces a convention that Provider structs do not have any exported fields, this is one part of helping to ensure secrets are not logged.
// json.Marshal for example will not include unexported values.
func TestProvidersHaveNoPublicAttributes(t *testing.T) {
	for _, tc := range Registry().Providers {
		t.Run(tc.DefaultID, func(t *testing.T) {
			v := reflect.ValueOf(tc.Provider)
			if v.Kind() == reflect.Ptr {
				if v.IsNil() {
					t.Fatal("unexpected nil provider in registry")
				}
				// dereference to get a value
				v = v.Elem()
			}
			// check for any exported fields
			for _, f := range reflect.VisibleFields(v.Type()) {
				assert.False(t, f.IsExported(), fmt.Sprintf("error in %s Provider struct. Field: '%s' should not be exported, change this to a lowercase name. By convention, all provider structs should not contain exported fields.", tc.DefaultID, f.Name))
			}
		})
	}
}
