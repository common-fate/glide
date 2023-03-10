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

func TestHealthcheckHandler(t *testing.T) {

	type testcase struct {
		name               string
		mockHealthcheckErr error
		want               string
		wantCode           int
	}

	testcases := []testcase{
		{
			name:     "ok",
			wantCode: http.StatusNoContent,
			want:     ``,
		},
		{
			name:               "error",
			mockHealthcheckErr: errors.New("an error"),
			wantCode:           http.StatusInternalServerError,
			want:               `{"error":"Internal Server Error"}`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mockHealthcheck := mocks.NewMockHealthcheckService(ctrl)
			mockHealthcheck.EXPECT().Check(gomock.Any()).Return(tc.mockHealthcheckErr).AnyTimes()
			a := API{HealthcheckService: mockHealthcheck}
			handler := newTestServer(t, &a)

			req, err := http.NewRequest("POST", "/api/v1/admin/healthcheck-handlers", nil)
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

func TestRegisterHandler(t *testing.T) {

	// test cases:
	// apio.DecodeJSONBody error ✅
	// RegisterHandler success ✅
	// RegisterHandler error == handlersvc.ErrHandlerIdAlreadyExists ✅
	// RegisterHandler error == anything else ✅

	type testcase struct {
		name                   string
		wantCode               int
		wantBody               string
		withRegisterResult     *handler.Handler
		giveBody               string
		mockRegisterHandlerErr error
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
			wantBody: `{"awsAccount":"123456789012","awsRegion":"ap-southeast-2","diagnostics":[],"functionArn":"arn:aws:lambda:ap-southeast-2:123456789012:function:handler","healthy":false,"id":"handler","runtime":"aws-lambda"}`,
			withRegisterResult: &handler.Handler{
				ID:          "handler",
				Runtime:     "aws-lambda",
				AWSAccount:  "123456789012",
				AWSRegion:   "ap-southeast-2",
				Healthy:     false,
				Diagnostics: []handler.Diagnostic{},
			},
			giveBody: `{"awsAccount":"123456789012","awsRegion":"ap-southeast-2","id":"handler","runtime":"aws-lambda"}`,
		},
		{
			name:               "invalid handlerID",
			wantCode:           http.StatusBadRequest,
			wantBody:           `{"error":"request body has an error: doesn't match the schema: Error at \"/id\": string doesn't match the regular expression \"^[-a-zA-Z0-9]*$\""}`,
			withRegisterResult: nil,
			giveBody:           `{"awsAccount":"123456789012","awsRegion":"ap-southeast-2","id":"handler with space","runtime":"aws-lambda"}`,
		},
		{
			name:                   "error == handlersvc.ErrHandlerIdAlreadyExists",
			mockRegisterHandlerErr: handlersvc.ErrHandlerIdAlreadyExists,
			wantCode:               http.StatusBadRequest,
			giveBody:               `{"awsAccount":"123456789012","awsRegion":"ap-southeast-2","id":"test","runtime":"aws-lambda"}`,
			wantBody:               `{"error":"handler id already exists"}`,
		},
		{
			name:                   "error == anything else",
			mockRegisterHandlerErr: errors.New("misc deployment svc error"),
			wantCode:               http.StatusInternalServerError,
			giveBody:               `{"awsAccount":"123456789012","awsRegion":"ap-southeast-2","id":"test","runtime":"aws-lambda"}`,
			wantBody:               `{"error":"Internal Server Error"}`,
		},
		{
			name:     "aws account validation too short",
			wantCode: http.StatusBadRequest,
			giveBody: `{"awsAccount":"123456789","awsRegion":"ap-southeast-2","id":"test","runtime":"aws-lambda"}`,
			wantBody: `{"error":"request body has an error: doesn't match the schema: Error at \"/awsAccount\": string doesn't match the regular expression \"^[0-9]{12}\""}`,
		},
		{
			name:     "aws account validation bad characters",
			wantCode: http.StatusBadRequest,
			giveBody: `{"awsAccount":"123456789abc","awsRegion":"ap-southeast-2","id":"test","runtime":"aws-lambda"}`,
			wantBody: `{"error":"request body has an error: doesn't match the schema: Error at \"/awsAccount\": string doesn't match the regular expression \"^[0-9]{12}\""}`,
		},
		{
			name:     "aws region validation",
			wantCode: http.StatusBadRequest,
			giveBody: `{"awsAccount":"123456789012","awsRegion":"ap-wrong-2","id":"test","runtime":"aws-lambda"}`,
			wantBody: `{"error":"request body has an error: doesn't match the schema: Error at \"/awsRegion\": string doesn't match the regular expression \"^(us(-gov)?|ap|ca|cn|eu|sa)-(central|(north|south)?(east|west)?)-\\d$\""}`,
		},
	}

	for _, tc := range testcases {

		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			// t.Parallel()

			ctrl := gomock.NewController(t)

			mockDeployment := mocks.NewMockHandlerService(ctrl)
			mockDeployment.EXPECT().RegisterHandler(gomock.Any(), gomock.Any()).Return(tc.withRegisterResult, tc.mockRegisterHandlerErr).AnyTimes()
			a := API{
				HandlerService: mockDeployment,
			}
			handler := newTestServer(t, &a)

			req, err := http.NewRequest(
				"POST",
				"/api/v1/admin/handlers",
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

func TestListHandlers(t *testing.T) {

	// test cases to handle
	// a.DB.Query(ctx, &q) error =  misc ✅
	// a.DB.Query(ctx, &q) error =  ddb.ErrNoItems ✅
	// a.DB.Query(ctx, &q) valid = 200 ✅

	type testcase struct {
		name        string
		handlers    []handler.Handler
		want        string
		mockListErr error
		wantCode    int
	}

	testcases := []testcase{
		{
			name:     "ok",
			wantCode: http.StatusOK,
			handlers: []handler.Handler{
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
			want: `{"next":"","res":[{"awsAccount":"string","awsRegion":"","diagnostics":[],"functionArn":"arn:aws:lambda::string:function:dep1","healthy":false,"id":"dep1","runtime":"string"},{"awsAccount":"string","awsRegion":"","diagnostics":[],"functionArn":"arn:aws:lambda::string:function:dep2","healthy":true,"id":"dep2","runtime":"string"}]}`,
		},
		{
			name:     "no handlers returns an empty list not an error",
			wantCode: http.StatusOK,
			handlers: nil,
			want:     `{"next":"","res":[]}`,
		},
		{
			name:        "internal error",
			mockListErr: errors.New("internal error"),
			wantCode:    http.StatusInternalServerError,
			handlers:    nil,

			want: `{"error":"Internal Server Error"}`,
		},
	}

	for _, tc := range testcases {

		// assign tc to a new variable so that it is not overwritten in the loop
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			db := ddbmock.New(t)
			db.MockQueryWithErr(&storage.ListHandlers{Result: tc.handlers}, tc.mockListErr)

			a := API{DB: db}
			handler := newTestServer(t, &a)

			req, err := http.NewRequest("GET", "/api/v1/admin/handlers", nil)
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

func TestGetHandler(t *testing.T) {

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

			req, err := http.NewRequest("GET", "/api/v1/admin/handlers/123", nil)
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

func TestDeleteHandler(t *testing.T) {

	type testcase struct {
		name                   string
		mockGetHandlerResponse handler.Handler
		mockGetHandlerErr      error
		mockDeleteHandlerErr   error
		want                   string
		wantCode               int
	}

	testcases := []testcase{
		{
			name:                   "ok",
			wantCode:               http.StatusNoContent,
			mockGetHandlerResponse: handler.Handler{ID: "123"},
			want:                   ``,
		},
		{
			name:              "deployment not found",
			wantCode:          http.StatusNotFound,
			mockGetHandlerErr: ddb.ErrNoItems,
			want:              `{"error":"item query returned no items"}`,
		},
		{
			name:              "internal error",
			wantCode:          http.StatusInternalServerError,
			mockGetHandlerErr: errors.New("internal error"),
			want:              `{"error":"Internal Server Error"}`,
		},
		{
			name:                   "internal error from delete",
			wantCode:               http.StatusInternalServerError,
			mockGetHandlerResponse: handler.Handler{ID: "123"},
			mockDeleteHandlerErr:   errors.New("some error"),
			want:                   `{"error":"Internal Server Error"}`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			db := ddbmock.New(t)
			db.MockQueryWithErr(&storage.GetHandler{Result: &tc.mockGetHandlerResponse}, tc.mockGetHandlerErr)
			ctrl := gomock.NewController(t)

			mockHandler := mocks.NewMockHandlerService(ctrl)
			mockHandler.EXPECT().DeleteHandler(gomock.Any(), &tc.mockGetHandlerResponse).Return(tc.mockDeleteHandlerErr).AnyTimes()
			a := API{DB: db, HandlerService: mockHandler}
			handler := newTestServer(t, &a)

			req, err := http.NewRequest("DELETE", "/api/v1/admin/handlers/123", nil)
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
