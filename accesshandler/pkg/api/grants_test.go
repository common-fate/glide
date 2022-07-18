package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/iso8601"

	"github.com/stretchr/testify/assert"
)

func TestPostGrants(t *testing.T) {
	type testcase struct {
		name     string
		body     string
		wantCode int
		wantErr  string
	}

	TenAM := time.Date(2022, 1, 1, 10, 0, 0, 0, time.UTC)

	TenAMISO8601 := iso8601.New(TenAM)
	TenThirtyAMISO8601 := iso8601.New(time.Date(2022, 1, 1, 10, 30, 0, 0, time.UTC))

	clk := clock.NewMock()
	clk.Set(TenAM)

	testcases := []testcase{
		{name: "ok", body: fmt.Sprintf(`{"id":"abcd","subject":"chris@commonfate.io","provider":"okta","with":{"group":"Admins"},"start":"%s","end":"%s"}`, TenAMISO8601, TenThirtyAMISO8601), wantCode: http.StatusCreated},
		{name: "invalid grant time: start after end", body: fmt.Sprintf(`{"id":"abcd","subject":"chris@commonfate.io","provider":"okta","with":{"group":"Admins"},"start":"%v","end":"%v"}`, TenThirtyAMISO8601, TenAMISO8601), wantCode: http.StatusBadRequest, wantErr: "grant start time must be earlier than end time"},
		{name: "invalid grant time: end before now", body: `{"id":"abcd","subject":"chris@commonfate.io","provider":"okta","with":{"group":"Admins"},"start": "2000-06-13T14:13:42.905Z","end": "2000-06-14T14:13:42.905Z"}`, wantCode: http.StatusBadRequest, wantErr: "grant finish time is in the past"},
		{name: "invalid body", body: `{"invalid": true}`, wantCode: http.StatusBadRequest, wantErr: "request body has an error: doesn't match the schema: Error at \"/subject\": property \"subject\" is missing"},
		{name: "missing subject", body: `{"id":"abcd","provider":"okta","with":{"group":"Admins"},"start":"2022-06-13T14:13:42.905Z","end":"2022-06-13T14:13:42.905Z"}`, wantCode: http.StatusBadRequest, wantErr: `request body has an error: doesn't match the schema: Error at "/subject": property "subject" is missing`},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			handler := newTestServer(t, withClock(clk))

			req, err := http.NewRequest("POST", "/api/v1/grants", strings.NewReader(tc.body))
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

func TestRevokeGrant(t *testing.T) {
	type testcase struct {
		name           string
		revokeBody     string
		body           string
		wantCode       int
		wantCodeRevoke int
		wantErr        string
	}

	TenAM := time.Date(2022, 1, 1, 10, 0, 0, 0, time.UTC)

	TenAMISO8601 := iso8601.New(TenAM)
	TenThirtyAMISO8601 := iso8601.New(time.Date(2022, 1, 1, 10, 30, 0, 0, time.UTC))

	clk := clock.NewMock()
	clk.Set(TenAM)

	testcases := []testcase{
		{name: "create grant and revoke ok", revokeBody: `{"revokerId":"1234"}`, body: fmt.Sprintf(`{"id":"abcd","subject":"chris@commonfate.io","provider":"okta","with":{"group":"Admins"},"start":"%s","end":"%s"}`, TenAMISO8601, TenThirtyAMISO8601), wantCode: http.StatusCreated, wantCodeRevoke: http.StatusOK},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			handler := newTestServer(t, withClock(clk))

			//create grant
			req, err := http.NewRequest("POST", "/api/v1/grants", strings.NewReader(tc.body))
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

			//revoke grant
			req, err = http.NewRequest("POST", "/api/v1/grants/123/revoke", strings.NewReader(tc.revokeBody))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Add("Content-Type", "application/json")

			rr = httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.wantCodeRevoke, rr.Code)

			_ = json.NewDecoder(rr.Body).Decode(&apiErr)
			assert.Equal(t, tc.wantErr, apiErr.Error)

		})
	}
}
