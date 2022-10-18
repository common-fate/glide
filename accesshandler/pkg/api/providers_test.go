package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/config"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers/testgroups"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/types"

	"github.com/stretchr/testify/assert"
)

func TestGetProvider(t *testing.T) {
	type testcase struct {
		name           string
		giveProviderId string
		wantCode       int
		wantErr        string
	}

	notFoundErr := &providers.ProviderNotFoundError{Provider: "badid"}

	testcases := []testcase{
		{name: "ok", giveProviderId: "test", wantCode: http.StatusOK},
		{name: "not found", giveProviderId: "badid", wantCode: http.StatusNotFound, wantErr: notFoundErr.Error()},
	}
	config.ConfigureTestProviders([]config.Provider{
		{
			ID:       "test",
			Type:     "testgroups",
			Provider: &testgroups.Provider{},
		},
	})

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			handler := newTestServer(t)

			req, err := http.NewRequest("GET", "/api/v1/providers/"+tc.giveProviderId, nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Add("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.wantCode, rr.Code)
			var apiErr apio.ErrorResponse

			_ = json.NewDecoder(rr.Body).Decode(&apiErr)
			assert.Equal(t, tc.wantErr, apiErr.Error)
		})
	}
}
func TestListProviders(t *testing.T) {
	type testcase struct {
		name     string
		wantCode int
		wantBody []types.Provider
	}

	testcases := []testcase{
		{name: "ok", wantCode: http.StatusOK, wantBody: []types.Provider{{Id: "test", Type: "testgroups"}}},
	}
	config.ConfigureTestProviders([]config.Provider{
		{
			ID:       "test",
			Type:     "testgroups",
			Provider: &testgroups.Provider{},
		},
	})
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			handler := newTestServer(t)

			req, err := http.NewRequest("GET", "/api/v1/providers", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Add("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.wantCode, rr.Code)

			var got []types.Provider
			_ = json.NewDecoder(rr.Body).Decode(&got)

			assert.Equal(t, tc.wantBody, got)
		})
	}
}
func TestGetProviderArgs(t *testing.T) {
	type testcase struct {
		name           string
		giveProviderId string
		wantBody       *types.ArgSchema
		wantCode       int
		wantErr        string
	}
	tg := &testgroups.Provider{}
	config.ConfigureTestProviders([]config.Provider{
		{
			ID:       "test",
			Type:     "testgroups",
			Provider: tg,
		},
	})

	notFoundErr := &providers.ProviderNotFoundError{Provider: "badid"}

	schema := tg.ArgSchema().ToAPI()
	testcases := []testcase{
		{name: "ok", giveProviderId: "test", wantCode: http.StatusOK, wantBody: &schema},
		{name: "not found", giveProviderId: "badid", wantCode: http.StatusNotFound, wantErr: notFoundErr.Error()},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			handler := newTestServer(t)

			req, err := http.NewRequest("GET", "/api/v1/providers/"+tc.giveProviderId+"/args", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Add("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.wantCode, rr.Code)
			var apiErr apio.ErrorResponse

			body, err := io.ReadAll(rr.Body)
			if err != nil {
				t.Fatal(err)
			}
			err = json.Unmarshal(body, &apiErr)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tc.wantErr, apiErr.Error)

			if tc.wantBody != nil {
				out, err := json.Marshal(tc.wantBody)
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, string(out), string(body))
			}
		})
	}
}

func TestListProviderArgOptions(t *testing.T) {
	type testcase struct {
		name           string
		giveProviderId string
		giveArgId      string
		wantBody       types.ArgOptionsResponse
		wantCode       int
		wantErr        string
	}

	notFoundErr := &providers.ProviderNotFoundError{Provider: "badid"}

	invalidArgErr := &providers.InvalidArgumentError{Arg: "notexist"}

	options := []types.Option{{Label: "group1", Value: "group1"}}
	tg := &testgroups.Provider{
		Groups: []string{"group1"},
	}
	config.ConfigureTestProviders([]config.Provider{
		{
			ID:       "test",
			Type:     "testgroups",
			Provider: tg,
		},
	})
	testcases := []testcase{
		{name: "ok", giveProviderId: "test", giveArgId: "group", wantCode: http.StatusOK, wantBody: types.ArgOptionsResponse{Options: options}},
		{name: "provider not found", giveProviderId: "badid", giveArgId: "notexist", wantCode: http.StatusNotFound, wantErr: notFoundErr.Error()},
		{name: "arg not found", giveProviderId: "test", giveArgId: "notexist", wantCode: http.StatusNotFound, wantErr: invalidArgErr.Error()},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			handler := newTestServer(t)

			req, err := http.NewRequest("GET", "/api/v1/providers/"+tc.giveProviderId+"/args/"+tc.giveArgId+"/options", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Add("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.wantCode, rr.Code)
			var apiErr apio.ErrorResponse

			body, err := io.ReadAll(rr.Body)
			if err != nil {
				t.Fatal(err)
			}
			if tc.wantErr != "" {
				err = json.Unmarshal(body, &apiErr)
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tc.wantErr, apiErr.Error)
			}

			var got types.ArgOptionsResponse
			err = json.Unmarshal(body, &got)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tc.wantBody, got)
		})
	}
}
