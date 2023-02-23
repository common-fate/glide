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
			give: `{"ID": "test", "targetSchema": "v1.0.1"}`,
			mockCreate: &target.Group{
				ID:           "test",
				TargetSchema: target.GroupTargetSchema{From: "v1.0.1", Schema: providerregistrysdk.TargetMode_Schema{}},
			},
			wantCode: http.StatusCreated,

			wantBody: `{"icon":"","id":"test","targetDeployments":null,"targetSchema":{"From":"v1.0.1","Schema":{}}}`,
		},
		{
			name:          "id already exists",
			give:          `{"ID": "test", "targetSchema": "v1.0.1"}`,
			mockCreateErr: targetsvc.ErrTargetGroupIdAlreadyExists,
			wantCode:      http.StatusConflict,
			wantBody:      `{"error":"target group id already exists"}`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mocks.NewMockTargetService(ctrl)
			m.EXPECT().CreateGroup(gomock.Any(), gomock.Any()).Return(tc.mockCreate, tc.mockCreateErr)

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
					ID:           "tg1",
					TargetSchema: target.GroupTargetSchema{From: "test", Schema: providerregistrysdk.TargetMode_Schema{AdditionalProperties: map[string]providerregistrysdk.TargetArgument{}}},
					Icon:         "test",
				},
				{
					ID:           "tg2",
					TargetSchema: target.GroupTargetSchema{From: "test", Schema: providerregistrysdk.TargetMode_Schema{AdditionalProperties: map[string]providerregistrysdk.TargetArgument{}}},
					Icon:         "test",
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
			t.Parallel()
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
			t.Parallel()
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
			want:                                 `{"diagnostics":[],"handlerId":"123","mode":"Default","priority":100,"targetGroupId":"123","valid":false}`,
			deploymentId:                         "abc",
			mockCreate:                           &target.Route{Group: "123", Handler: "123", Kind: "Default", Priority: 100},
			give:                                 `{"deploymentId": "abc", "priority": 100,"force":false}`,
		},

		{
			name:         "priority cannot be out of range",
			wantCode:     http.StatusBadRequest,
			want:         `{"error":"request body has an error: doesn't match the schema: Error at \"/priority\": number must be at most 999"}`,
			deploymentId: "abc",
			give:         `{"deploymentId": "abc", "priority": 1000,"force":false}`,
		},
		{
			name:         "priority cannot be under range",
			wantCode:     http.StatusBadRequest,
			want:         `{"error":"request body has an error: doesn't match the schema: Error at \"/priority\": number must be at least 0"}`,
			deploymentId: "abc",
			give:         `{"deploymentId": "abc", "priority": -1,"force":false}`,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
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
		mockCreateErr                        error
	}

	testcases := []testcase{
		{
			name:                                 "ok",
			wantCode:                             http.StatusOK,
			mockGetTargetGroupResponse:           target.Group{ID: "123"},
			mockGetTargetGroupDeploymentResponse: handler.Handler{ID: "abc"},
			want:                                 `null`,
			deploymentId:                         "abc",
			mockCreate:                           &target.Group{ID: "123"},
		},
		{
			name:                                 "target group err, error case",
			wantCode:                             http.StatusInternalServerError,
			mockCreateErr:                        errors.New("error case"),
			mockGetTargetGroupDeploymentResponse: handler.Handler{ID: "abc"},
			want:                                 `{"error":"Internal Server Error"}`,
			deploymentId:                         "abc",
			mockCreate:                           &target.Group{ID: "123"},
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			db := ddbmock.New(t)
			db.MockQueryWithErr(&storage.GetTargetGroup{Result: &tc.mockGetTargetGroupResponse}, tc.mockGetTargetGroupErr)
			db.MockQueryWithErr(&storage.GetHandler{Result: &tc.mockGetTargetGroupDeploymentResponse}, tc.mockGetTargetGroupErr)

			a := API{DB: db}
			handler := newTestServer(t, &a)

			req, err := http.NewRequest("POST", fmt.Sprintf("/api/v1/admin/target-groups/123/unlink?deploymentId=%s", tc.deploymentId), strings.NewReader(""))
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
