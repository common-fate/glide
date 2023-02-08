package api

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/common-fate/pkg/api/mocks"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/targetgroup"
	"github.com/common-fate/ddb"
	"github.com/common-fate/ddb/ddbmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// func TestListTargetGroupDeployments(t *testing.T) {
func TestCreateTargetGroupDeployments(t *testing.T) {
	type testcase struct {
		name           string
		mockCancelErr  error
		wantCode       int
		wantBody       string
		withCreatedDep *targetgroup.Deployment
	}

	testcases := []testcase{
		{
			name:     "create.success.201",
			wantCode: http.StatusCreated,
			wantBody: `{"id":"123456789012","functionArn":"string","runtime":"string","awsAccount":"string","healthy":false,"diagnostics":[],"activeConfig":{},"provider":{"publisher":"","name":"","version":""}}`,
			withCreatedDep: &targetgroup.Deployment{
				ID:           "123456789012",
				FunctionARN:  "string",
				Runtime:      "string",
				AWSAccount:   "string",
				Healthy:      false,
				Diagnostics:  []targetgroup.Diagnostic{},
				ActiveConfig: map[string]targetgroup.Config{},
				Provider:     targetgroup.Provider{},
			},
		},
		{
			name:          "internal.error.500",
			mockCancelErr: apio.NewRequestError(errors.New("oh no"), http.StatusInternalServerError),
			wantCode:      http.StatusInternalServerError,
			wantBody:      `{"error":"oh no"}`,
		},
	}

	for _, tc := range testcases {

		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)

			mockDeployment := mocks.NewMockTargetGroupDeploymentService(ctrl)
			mockDeployment.EXPECT().CreateTargetGroupDeployment(gomock.Any(), gomock.Any()).Return(tc.withCreatedDep, tc.mockCancelErr).AnyTimes()
			a := API{
				TargetGroupDeploymentService: mockDeployment,
			}
			handler := newTestServer(t, &a)

			req, err := http.NewRequest(
				"POST",
				"/api/v1/target-group-deployments",
				// we pass a nullish body here (this is not used due to the service mocks)
				strings.NewReader("{}"),
			)

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

func TestListTargetGroupDeployments(t *testing.T) {

	type testcase struct {
		name                   string
		targetGroupDeployments []targetgroup.Deployment
		want                   string
		mockListErr            error
		wantCode               int
	}

	testcases := []testcase{
		{
			name:     "ok",
			wantCode: http.StatusOK,
			targetGroupDeployments: []targetgroup.Deployment{
				{
					ID:           "dep1",
					FunctionARN:  "string",
					Runtime:      "string",
					AWSAccount:   "string",
					Healthy:      false,
					Diagnostics:  []targetgroup.Diagnostic{},
					ActiveConfig: map[string]targetgroup.Config{},
					Provider:     targetgroup.Provider{},
				},
				{
					ID:           "dep2",
					FunctionARN:  "string",
					Runtime:      "string",
					AWSAccount:   "string",
					Healthy:      true,
					Diagnostics:  []targetgroup.Diagnostic{},
					ActiveConfig: map[string]targetgroup.Config{},
					Provider:     targetgroup.Provider{},
				},
			},
			want: `{"next":"","res":[{"activeConfig":{"test":{"type":"test","value":{}}},"awsAccount":"string","diagnostics":[],"functionArn":"string","healthy":false,"id":"dep1","provider":{"name":"","publisher":"","version":""}},{"activeConfig":{"test":{"type":"test","value":{}}},"awsAccount":"string","diagnostics":[],"functionArn":"string","healthy":true,"id":"dep2","provider":{"name":"","publisher":"","version":""}}]}`,
		},
		{
			name:                   "no target groups returns an empty list not an error",
			mockListErr:            ddb.ErrNoItems,
			wantCode:               http.StatusOK,
			targetGroupDeployments: nil,

			want: `{"next":"","res":[]}`,
		},
		{
			name:                   "internal error",
			mockListErr:            errors.New("internal error"),
			wantCode:               http.StatusInternalServerError,
			targetGroupDeployments: nil,

			want: `{"error":"Internal Server Error"}`,
		},
	}

	for _, tc := range testcases {

		// assign tc to a new variable so that it is not overwritten in the loop
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			db := ddbmock.New(t)
			db.MockQueryWithErr(&storage.ListTargetGroupDeployments{Result: tc.targetGroupDeployments}, tc.mockListErr)

			a := API{DB: db}
			handler := newTestServer(t, &a)

			req, err := http.NewRequest("GET", "/api/v1/target-group-deployments", nil)
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

func TestGetTargetGroupDeployment(t *testing.T) {

	type testcase struct {
		name                          string
		mockGetTargetGroupDepResponse targetgroup.Deployment
		mockGetTargetGroupDepErr      error
		want                          string
		wantCode                      int
	}

	testcases := []testcase{
		{
			name:                          "ok",
			wantCode:                      http.StatusOK,
			mockGetTargetGroupDepResponse: targetgroup.Deployment{ID: "123"},
			want:                          `{"icon":"","id":"123","targetDeployments":null,"targetSchema":{"From":"","Schema":{}}}`,
		},
		{
			name:                     "deployment not found",
			wantCode:                 http.StatusNotFound,
			mockGetTargetGroupDepErr: ddb.ErrNoItems,
			want:                     `{"error":"item query returned no items"}`,
		},
		{
			name:                     "internal error",
			wantCode:                 http.StatusInternalServerError,
			mockGetTargetGroupDepErr: errors.New("internal error"),
			want:                     `{"error":"Internal Server Error"}`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			db := ddbmock.New(t)
			db.MockQueryWithErr(&storage.GetTargetGroupDeployment{Result: tc.mockGetTargetGroupDepResponse}, tc.mockGetTargetGroupDepErr)

			a := API{DB: db}
			handler := newTestServer(t, &a)

			req, err := http.NewRequest("GET", "/api/v1/target-group-deployments/123", nil)
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
