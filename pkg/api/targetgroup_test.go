package api

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/common-fate/common-fate/pkg/api/mocks"
	"github.com/common-fate/common-fate/pkg/handler"
	"github.com/common-fate/common-fate/pkg/service/targetsvc"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/target"

	"github.com/common-fate/ddb"
	"github.com/common-fate/ddb/ddbmock"
	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateTargetGroup(t *testing.T) {
	type testcase struct {
		name          string
		give          string
		mockCreate    *target.Group
		mockCreateErr error

		wantCode int
		wantBody string
	}

	testcases := []testcase{
		{
			name: "ok",
			give: `{"id": "test", "from": {"publisher": "common-fate", "name": "test", "version": "v1", "kind": "Kind"}}`,
			mockCreate: &target.Group{
				ID: "test",
				From: target.From{
					Name: "test",
				},
				Schema: providerregistrysdk.Target{},
			},
			wantCode: http.StatusCreated,

			wantBody: `{"createdAt":"0001-01-01T00:00:00Z","from":{"kind":"","name":"test","publisher":"","version":""},"icon":"","id":"test","schema":{},"updatedAt":"0001-01-01T00:00:00Z"}`,
		},
		{
			name:     "invalid-target-id",
			give:     `{"id": "target id with space", "from": {"publisher": "common-fate", "name": "test", "version": "v1", "kind": "Kind"}}`,
			wantCode: http.StatusBadRequest,
			wantBody: `{"error":"request body has an error: doesn't match the schema: Error at \"/id\": string doesn't match the regular expression \"^[-a-zA-Z0-9]*$\""}`,
		},
		{
			name:     "maximum length exceeded for target id",
			give:     `{"id": "target-id-max-length-test-target-id-max-length-test-target-id-max-length-test-target-id-max-length-test-target-id-max-length-test-target-id-max-length-test-", "from": {"publisher": "common-fate", "name": "test", "version": "v1", "kind": "Kind"}}`,
			wantCode: http.StatusBadRequest,
			wantBody: `{"error":"request body has an error: doesn't match the schema: Error at \"/id\": maximum string length is 64"}`,
		},
		{
			name:          "id already exists",
			give:          `{"id": "test", "from": {"publisher": "common-fate", "name": "test", "version": "v1", "kind": "Kind"}}`,
			mockCreateErr: targetsvc.ErrTargetGroupIdAlreadyExists,
			wantCode:      http.StatusConflict,
			wantBody:      `{"error":"target group id already exists"}`,
		},
		{
			name:          "provider not found in registry",
			give:          `{"id": "test", "from": {"publisher": "common-fate", "name": "test", "version": "v1", "kind": "Kind"}}`,
			mockCreateErr: targetsvc.ErrProviderNotFoundInRegistry,
			wantCode:      http.StatusNotFound,
			wantBody:      `{"error":"provider not found in registry"}`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			//
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mocks.NewMockTargetService(ctrl)

			if tc.mockCreate != nil || tc.mockCreateErr != nil {
				m.EXPECT().CreateGroup(gomock.Any(), gomock.Any()).Return(tc.mockCreate, tc.mockCreateErr)
			}

			a := API{TargetService: m}
			handler := newTestServer(t, &a)

			req, err := http.NewRequest("POST", "/api/v1/admin/target-groups", strings.NewReader(tc.give))
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

func TestListTargetGroup(t *testing.T) {
	type testcase struct {
		name string

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

			want: `{"targetGroups":[{"icon":"test","id":"tg1","targetDeployments":[{"Diagnostics":null,"Id":"reg1","Priority":0,"Valid":false}],"targetSchema":{"From":"test","Schema":{}}},{"icon":"test","id":"tg2","targetDeployments":[{"Diagnostics":null,"Id":"reg1","Priority":0,"Valid":false}],"targetSchema":{"From":"test","Schema":{}}}]}`,
		},
		{
			name:         "no target groups returns an empty list not an error",
			mockListErr:  ddb.ErrNoItems,
			wantCode:     http.StatusOK,
			targetgroups: nil,

			want: `{"targetGroups":[]}`,
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

			req, err := http.NewRequest("GET", "/api/v1/admin/target-groups", nil)
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

func TestGetTargetGroup(t *testing.T) {
	type testcase struct {
		name                       string
		mockGetTargetGroupResponse target.Group
		mockGetTargetGroupErr      error
		want                       string
		wantCode                   int
	}

	testcases := []testcase{
		{
			name:                       "ok",
			wantCode:                   http.StatusOK,
			mockGetTargetGroupResponse: target.Group{ID: "123"},
			want:                       `{"icon":"","id":"123","targetDeployments":null,"targetSchema":{"From":"","Schema":{}}}`,
		},
		{
			name:                  "group not found",
			wantCode:              http.StatusNotFound,
			mockGetTargetGroupErr: ddb.ErrNoItems,

			want: `{"error":"item query returned no items"}`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {

			db := ddbmock.New(t)
			db.MockQueryWithErr(&storage.GetTargetGroup{Result: &tc.mockGetTargetGroupResponse}, tc.mockGetTargetGroupErr)

			a := API{DB: db}
			handler := newTestServer(t, &a)

			req, err := http.NewRequest("GET", "/api/v1/admin/target-groups/123", nil)
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

func TestTargetGroupLink(t *testing.T) {
	type testcase struct {
		name                                 string
		mockGetTargetGroupResponse           target.Group
		mockGetTargetGroupDeploymentResponse handler.Handler
		mockGetTargetGroupErr                error
		want                                 string
		wantCode                             int
		mockCreate                           *target.Route
		deploymentId                         string
		mockCreateErr                        error
		give                                 string
	}

	testcases := []testcase{
		{
			name:                                 "ok",
			wantCode:                             http.StatusOK,
			mockGetTargetGroupResponse:           target.Group{ID: "123"},
			mockGetTargetGroupDeploymentResponse: handler.Handler{ID: "abc"},
			want:                                 `{"diagnostics":[],"handlerId":"123","kind":"Default","priority":100,"targetGroupId":"123","valid":false}`,
			deploymentId:                         "abc",
			mockCreate:                           &target.Route{Group: "123", Handler: "123", Kind: "Default", Priority: 100},
			give:                                 `{"deploymentId": "abc", "priority": 100, "kind":"Default"}`,
		},

		{
			name:         "priority cannot be out of range",
			wantCode:     http.StatusBadRequest,
			want:         `{"error":"request body has an error: doesn't match the schema: Error at \"/priority\": number must be at most 999"}`,
			deploymentId: "abc",
			give:         `{"deploymentId": "abc", "priority": 1000, "kind":"Default"}`,
		},
		{
			name:         "priority cannot be under range",
			wantCode:     http.StatusBadRequest,
			want:         `{"error":"request body has an error: doesn't match the schema: Error at \"/priority\": number must be at least 0"}`,
			deploymentId: "abc",
			give:         `{"deploymentId": "abc", "priority": -1, "kind":"Default"}`,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			db := ddbmock.New(t)
			db.MockQueryWithErr(&storage.GetTargetGroup{Result: &tc.mockGetTargetGroupResponse}, tc.mockGetTargetGroupErr)
			db.MockQueryWithErr(&storage.GetHandler{Result: &tc.mockGetTargetGroupDeploymentResponse}, tc.mockGetTargetGroupErr)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := mocks.NewMockTargetService(ctrl)

			m.EXPECT().CreateRoute(gomock.Any(), gomock.Any(), gomock.Any()).Return(tc.mockCreate, tc.mockCreateErr).AnyTimes()

			a := API{DB: db, TargetService: m}
			handler := newTestServer(t, &a)

			req, err := http.NewRequest("POST", "/api/v1/admin/target-groups/123/link", strings.NewReader(tc.give))
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

func TestRemoveTargetGroupLink(t *testing.T) {
	type testcase struct {
		name                                 string
		mockGetTargetGroupResponse           target.Group
		mockGetTargetGroupDeploymentResponse handler.Handler
		mockGetTargetGroupErr                error
		want                                 string
		wantCode                             int
		mockCreate                           *target.Group
		deploymentId                         string
		kind                                 string
	}

	testcases := []testcase{
		{
			name:                                 "ok",
			wantCode:                             http.StatusOK,
			mockGetTargetGroupResponse:           target.Group{ID: "123"},
			mockGetTargetGroupDeploymentResponse: handler.Handler{ID: "abc"},
			want:                                 `null`,
			deploymentId:                         "abc",
			kind:                                 "Default",
			mockCreate:                           &target.Group{ID: "123"},
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			db := ddbmock.New(t)
			db.MockQueryWithErr(&storage.GetTargetGroup{Result: &tc.mockGetTargetGroupResponse}, tc.mockGetTargetGroupErr)
			db.MockQueryWithErr(&storage.GetHandler{Result: &tc.mockGetTargetGroupDeploymentResponse}, tc.mockGetTargetGroupErr)

			a := API{DB: db}
			handler := newTestServer(t, &a)

			req, err := http.NewRequest("POST", fmt.Sprintf("/api/v1/admin/target-groups/123/unlink?deploymentId=%s&kind=%s", tc.deploymentId, tc.kind), strings.NewReader(""))
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

func TestDeleteTargetGroup(t *testing.T) {
	type testcase struct {
		name                     string
		mockGetTargetGroup       *target.Group
		mockGetTargetGroupErr    error
		mockDeleteTargetGroupErr error
		want                     string
		wantCode                 int
	}

	testcases := []testcase{
		{
			name:     "ok",
			wantCode: http.StatusNoContent,
			want:     ``,
		},
		{
			name:                  "not found",
			wantCode:              http.StatusNotFound,
			mockGetTargetGroupErr: ddb.ErrNoItems,
			want:                  `{"error":"item query returned no items"}`,
		},
		{
			name:                     "internal error",
			wantCode:                 http.StatusInternalServerError,
			mockDeleteTargetGroupErr: errors.New("some error"),
			want:                     `{"error":"Internal Server Error"}`,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			db := ddbmock.New(t)
			db.MockQueryWithErr(&storage.GetTargetGroup{Result: tc.mockGetTargetGroup}, tc.mockGetTargetGroupErr)
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mocks.NewMockTargetService(ctrl)
			m.EXPECT().DeleteGroup(gomock.Any(), gomock.Any()).Return(tc.mockDeleteTargetGroupErr).AnyTimes()
			a := API{DB: db, TargetService: m}
			handler := newTestServer(t, &a)

			req, err := http.NewRequest("DELETE", "/api/v1/admin/target-groups/123", strings.NewReader(""))
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

func TestListTargetGroupRoutes(t *testing.T) {
	type testcase struct {
		name        string
		routes      []target.Route
		want        string
		mockListErr error
		wantCode    int
	}

	testcases := []testcase{
		{
			name:     "ok",
			wantCode: http.StatusOK,
			routes: []target.Route{
				{
					Group:       "abc",
					Handler:     "123",
					Kind:        "test",
					Priority:    999,
					Valid:       true,
					Diagnostics: []target.Diagnostic{},
				},
				{
					Group:       "abc",
					Handler:     "123",
					Kind:        "test",
					Priority:    999,
					Valid:       true,
					Diagnostics: []target.Diagnostic{},
				},
			},

			want: `{"routes":[{"diagnostics":[],"handlerId":"123","kind":"test","priority":999,"targetGroupId":"abc","valid":true},{"diagnostics":[],"handlerId":"123","kind":"test","priority":999,"targetGroupId":"abc","valid":true}]}`,
		},
		{
			name:        "no routes returns an empty list not an error",
			mockListErr: ddb.ErrNoItems,
			wantCode:    http.StatusOK,
			routes:      nil,

			want: `{"routes":[]}`,
		},
		{
			name:        "internal error",
			mockListErr: errors.New("internal error"),
			wantCode:    http.StatusInternalServerError,
			routes:      nil,

			want: `{"error":"Internal Server Error"}`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {

			db := ddbmock.New(t)
			db.MockQueryWithErr(&storage.ListTargetRoutesForGroup{Result: tc.routes}, tc.mockListErr)

			a := API{DB: db}
			handler := newTestServer(t, &a)

			req, err := http.NewRequest("GET", "/api/v1/admin/target-groups/abc/routes", nil)
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
