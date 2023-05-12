package api

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/common-fate/common-fate/pkg/cache"
	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/target"
	"github.com/common-fate/ddb"
	"github.com/common-fate/ddb/ddbmock"
	"github.com/stretchr/testify/assert"
)

func TestListEntitlements(t *testing.T) {
	type testcase struct {
		name         string
		targetgroups []target.Group
		want         string
		mockListErr  error
		wantCode     int
	}

	testcases := []testcase{
		{
			name:     "ok",
			wantCode: http.StatusOK,
			targetgroups: []target.Group{
				{
					ID: "tg1",
					From: target.From{
						Publisher: "common-fate",
						Name:      "test",
						Version:   "v1",
						Kind:      "Kind",
					},
					Icon: "test",
				},
				{
					ID: "tg2",
					From: target.From{
						Publisher: "common-fate",
						Name:      "second",
						Version:   "v2",
						Kind:      "Kind",
					},
					Icon: "test",
				},
			},

			want: `{"entitlements":[{"icon":"test","kind":"Kind","name":"test","publisher":"common-fate"},{"icon":"test","kind":"Kind","name":"second","publisher":"common-fate"}]}`,
		},
		{
			name:         "no entitlements returns an empty list not an error",
			mockListErr:  ddb.ErrNoItems,
			wantCode:     http.StatusOK,
			targetgroups: []target.Group{},

			want: `{"entitlements":[]}`,
		},
		{
			name:         "internal error",
			mockListErr:  errors.New("internal error"),
			wantCode:     http.StatusInternalServerError,
			targetgroups: nil,

			want: `{"error":"Internal Server Error"}`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {

			db := ddbmock.New(t)
			db.MockQueryWithErr(&storage.ListTargetGroups{Result: tc.targetgroups}, tc.mockListErr)

			a := API{DB: db}
			handler := newTestServer(t, &a)

			req, err := http.NewRequest("GET", "/api/v1/entitlements", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Add("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.wantCode, rr.Code)

			data, err := io.ReadAll(rr.Body)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tc.want, string(data))
		})
	}
}

var AccessRulesMap = make(map[string]cache.AccessRule)

// db: ListCachedTargetsForKind
// db: ListCachedTargets
// response: ListTargetResponse{}
func TestListEntitlementTargets(t *testing.T) {
	type testcase struct {
		name         string
		targets      []cache.Target
		want         string
		withTestUser *identity.User
		mockListErr  error
		wantCode     int
	}

	testcases := []testcase{
		{
			name:         "ok",
			withTestUser: &identity.User{Groups: []string{"testAdmin"}},
			targets: []cache.Target{
				{
					Fields: []cache.Field{
						{
							ID:         "id",
							FieldTitle: "account",
							ValueLabel: "account",
							Value:      "0123",
						},
					},
					AccessRules: map[string]cache.AccessRule{
						"foo": {
							MatchedTargetGroups: []string{"id"},
						},
					},
					IDPGroupsWithAccess: map[string]struct{}{"testAdmin": {}},
				},
			},
			want:        `{"targets":[{"fields":[{"fieldTitle":"account","id":"id","value":"0123","valueLabel":"account"}],"id":"###id#0123#","kind":{"icon":"","kind":"","name":"","publisher":""}}]}`,
			mockListErr: nil,
			wantCode:    http.StatusOK,
		},
		{
			name:        "no entitlements returns an empty list not an error",
			mockListErr: ddb.ErrNoItems,
			wantCode:    http.StatusOK,
			targets:     []cache.Target{},

			want: `{"targets":[]}`,
		},
		{
			name:        "internal error",
			mockListErr: errors.New("internal error"),
			wantCode:    http.StatusInternalServerError,
			targets:     nil,

			want: `{"error":"Internal Server Error"}`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {

			db := ddbmock.New(t)
			db.MockQueryWithErr(&storage.ListCachedTargets{Result: tc.targets}, tc.mockListErr)
			db.MockQueryWithErr(&storage.ListCachedTargetsForKind{Result: tc.targets}, tc.mockListErr)

			opts := []func(*testOptions){}
			if tc.withTestUser != nil {
				opts = append(opts, WithRequestUser(*tc.withTestUser))
			}

			a := API{DB: db}
			handler := newTestServer(t, &a, opts...)

			req, err := http.NewRequest("GET", "/api/v1/entitlements/targets", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Add("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.wantCode, rr.Code)

			data, err := io.ReadAll(rr.Body)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tc.want, string(data))
		})
	}
}
