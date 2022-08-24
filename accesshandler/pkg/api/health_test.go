package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealth(t *testing.T) {
	type testcase struct {
		name     string
		wantCode int
		wantBody string
	}

	testcases := []testcase{
		{name: "ok", wantCode: http.StatusOK, wantBody: `{"health":{"error":null,"healthy":true,"id":"okta"}}`},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			handler := newTestServer(t)

			req, err := http.NewRequest("GET", "/api/v1/health", nil)
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

			assert.Equal(t, tc.wantBody, string(data))
		})
	}
}
