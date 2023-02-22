package api

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/common-fate/common-fate/pkg/api/mocks"
	"github.com/common-fate/common-fate/pkg/handler"
	"github.com/common-fate/common-fate/pkg/service/handlersvc"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/ddb"
	"github.com/common-fate/ddb/ddbmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// func TestListTargetGroupDeployments(t *testing.T) {
func TestCreateTargetGroupDeployments(t *testing.T) {

	// test cases:
	// apio.DecodeJSONBody error ✅
	// CreateTargetGroupDeployment success ✅
	// CreateTargetGroupDeployment error == handlersvc.ErrTargetGroupDeploymentIdAlreadyExists ✅
	// CreateTargetGroupDeployment error == anything else ✅

	// items to mock: a.TargetGroupDeploymentService.CreateTargetGroupDeployment:res
	// items to mock: a.TargetGroupDeploymentService.CreateTargetGroupDeployment:err
	// items to test: request as types.CreateTargetGroupDeploymentRequest{}

	// we will need to mock a body like so:
	// giveBody: `{"createdAt":"0001-01-01T00:00:00Z","icon":"","id":"123","targetSchema":{"From":"","Schema":{}},"updatedAt":"0001-01-01T00:00:00Z"}`,
	// give that writing out a string like this is long and arduous, we will use a helper function that converts,
	// from types.CreateTargetGroupDeploymentRequest to a json object encoded as string
	// we can then parse giveBody as a types.CreateTargetGroupDeploymentRequest{} and compare it to the request
	// we will need to use a mockClock for consistent createdAt and updatedAt values

	type testcase struct {
		name           string
		wantCode       int
		wantBody       string
		withCreatedDep *handler.Handler
		giveBody       string
		// mockCreateTargetgroupDeployment    *handler.Handler
		mockCreateTargetgroupDeploymentErr error
	}

	testcases := []testcase{
		{
			name:     "apio.DecodeJSONBody error",
			wantCode: http.StatusBadRequest,
			wantBody: `{"error":"request body has an error: doesn't match the schema: Error at \"/id\": property \"id\" is missing"}`,
			giveBody: "{}",
		},
		{
			name:     "create.success.201",
			wantCode: http.StatusCreated,
			wantBody: `{"awsAccount":"string","awsRegion":"","diagnostics":[],"functionArn":"arn:aws:lambda::string:function:123456789012","healthy":false,"id":"123456789012"}`,
			withCreatedDep: &handler.Handler{
				ID:          "123456789012",
				Runtime:     "string",
				AWSAccount:  "string",
				Healthy:     false,
				Diagnostics: []handler.Diagnostic{},
			},
			giveBody: `{"awsAccount":"123456789012","awsRegion":"ap-southeast-2","id":"test","runtime":"aws-lambda"}`,
		},
		{
			name:                               "error == handlersvc.ErrTargetGroupDeploymentIdAlreadyExists",
			mockCreateTargetgroupDeploymentErr: handlersvc.ErrHandlerIdAlreadyExists,
			wantCode:                           http.StatusBadRequest,
			giveBody:                           `{"awsAccount":"123456789012","awsRegion":"ap-southeast-2","id":"test","runtime":"aws-lambda"}`,
			wantBody:                           `{"error":"target group deployment id already exists"}`,
		},
		{
			name:                               "error == anything else",
			mockCreateTargetgroupDeploymentErr: errors.New("misc deployment svc error"),
			wantCode:                           http.StatusInternalServerError,
			giveBody:                           `{"awsAccount":"123456789012","awsRegion":"ap-southeast-2","id":"test","runtime":"aws-lambda"}`,
			wantBody:                           `{"error":"Internal Server Error"}`,
		},
	}

	for _, tc := range testcases {

		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)

			mockDeployment := mocks.NewMockHandlerService(ctrl)
			mockDeployment.EXPECT().RegisterHandler(gomock.Any(), gomock.Any()).Return(tc.withCreatedDep, tc.mockCreateTargetgroupDeploymentErr).AnyTimes()
			a := API{
				// TargetGroupDeploymentService: mockDeployment,
			}
			handler := newTestServer(t, &a)

			req, err := http.NewRequest(
				"POST",
				"/api/v1/admin/target-group-deployments",
				strings.NewReader(tc.giveBody),
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

	// test cases to handle
	// a.DB.Query(ctx, &q) error =  misc ✅
	// a.DB.Query(ctx, &q) error =  ddb.ErrNoItems ✅
	// a.DB.Query(ctx, &q) valid = 200 ✅

	type testcase struct {
		name                   string
		targetGroupDeployments []handler.Handler
		want                   string
		mockListErr            error
		wantCode               int
	}

	testcases := []testcase{
		{
			name:     "ok",
			wantCode: http.StatusOK,
			targetGroupDeployments: []handler.Handler{
				{
					ID:          "dep1",
					Runtime:     "string",
					AWSAccount:  "string",
					Healthy:     false,
					Diagnostics: []handler.Diagnostic{},
				},
				{
					ID:          "dep2",
					Runtime:     "string",
					AWSAccount:  "string",
					Healthy:     true,
					Diagnostics: []handler.Diagnostic{},
				},
			},
			want: `{"next":"","res":[{"awsAccount":"string","awsRegion":"","diagnostics":[],"functionArn":"arn:aws:lambda::string:function:dep1","healthy":false,"id":"dep1"},{"awsAccount":"string","awsRegion":"","diagnostics":[],"functionArn":"arn:aws:lambda::string:function:dep2","healthy":true,"id":"dep2"}]}`,
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
			db.MockQueryWithErr(&storage.ListHandlers{Result: tc.targetGroupDeployments}, tc.mockListErr)

			a := API{DB: db}
			handler := newTestServer(t, &a)

			req, err := http.NewRequest("GET", "/api/v1/admin/target-group-deployments", nil)
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
		mockGetTargetGroupDepResponse handler.Handler
		mockGetTargetGroupDepErr      error
		want                          string
		wantCode                      int
	}

	testcases := []testcase{
		{
			name:                          "ok",
			wantCode:                      http.StatusOK,
			mockGetTargetGroupDepResponse: handler.Handler{ID: "123"},
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
			db.MockQueryWithErr(&storage.GetHandler{Result: &tc.mockGetTargetGroupDepResponse}, tc.mockGetTargetGroupDepErr)

			a := API{DB: db}
			handler := newTestServer(t, &a)

			req, err := http.NewRequest("GET", "/api/v1/admin/target-group-deployments/123", nil)
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
