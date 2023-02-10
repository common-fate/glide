package api

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/common-fate/common-fate/pkg/api/mocks"
	"github.com/common-fate/common-fate/pkg/service/targetgroupsvc"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/targetgroup"
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
		mockCreate    *targetgroup.TargetGroup
		mockCreateErr error

		wantCode int
		wantBody string
	}

	testcases := []testcase{
		{
			name: "ok",
			give: `{"ID": "test", "targetSchema": "v1.0.1"}`,
			mockCreate: &targetgroup.TargetGroup{
				ID:           "test",
				TargetSchema: targetgroup.GroupTargetSchema{From: "v1.0.1", Schema: providerregistrysdk.TargetSchema{}},
			},
			wantCode: http.StatusCreated,

			wantBody: `{"icon":"","id":"test","targetDeployments":null,"targetSchema":{"From":"v1.0.1","Schema":{}}}`,
		},
		{
			name:          "id already exists",
			give:          `{"ID": "test", "targetSchema": "v1.0.1"}`,
			mockCreateErr: targetgroupsvc.ErrTargetGroupIdAlreadyExists,
			wantCode:      http.StatusBadRequest,
			wantBody:      `{"error":"target group id already exists"}`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mocks.NewMockTargetGroupService(ctrl)
			m.EXPECT().CreateTargetGroup(gomock.Any(), gomock.Any()).Return(tc.mockCreate, tc.mockCreateErr)

			a := API{TargetGroupService: m}
			handler := newTestServer(t, &a)

			req, err := http.NewRequest("POST", "/api/v1/target-groups", strings.NewReader(tc.give))
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

		targetgroups []targetgroup.TargetGroup
		want         string
		mockListErr  error
		wantCode     int
	}

	testcases := []testcase{
		{
			name:     "ok",
			wantCode: http.StatusOK,
			targetgroups: []targetgroup.TargetGroup{
				{
					ID:           "tg1",
					TargetSchema: targetgroup.GroupTargetSchema{From: "test", Schema: providerregistrysdk.TargetSchema{AdditionalProperties: map[string]providerregistrysdk.TargetArgument{}}},
					Icon:         "test",
				},
				{
					ID:           "tg2",
					TargetSchema: targetgroup.GroupTargetSchema{From: "test", Schema: providerregistrysdk.TargetSchema{AdditionalProperties: map[string]providerregistrysdk.TargetArgument{}}},
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

			req, err := http.NewRequest("GET", "/api/v1/target-groups", nil)
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
		mockGetTargetGroupResponse targetgroup.TargetGroup
		mockGetTargetGroupErr      error
		want                       string
		wantCode                   int
	}

	testcases := []testcase{
		{
			name:                       "ok",
			wantCode:                   http.StatusOK,
			mockGetTargetGroupResponse: targetgroup.TargetGroup{ID: "123"},
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
			db.MockQueryWithErr(&storage.GetTargetGroup{Result: tc.mockGetTargetGroupResponse}, tc.mockGetTargetGroupErr)

			a := API{DB: db}
			handler := newTestServer(t, &a)

			req, err := http.NewRequest("GET", "/api/v1/target-groups/123", nil)
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
		mockGetTargetGroupResponse           targetgroup.TargetGroup
		mockGetTargetGroupDeploymentResponse targetgroup.Deployment
		mockGetTargetGroupErr                error
		want                                 string
		wantCode                             int
		mockCreate                           *targetgroup.TargetGroup
		deploymentId                         string
		mockCreateErr                        error
		give                                 string
	}

	testcases := []testcase{
		{
			name:                                 "ok",
			wantCode:                             http.StatusOK,
			mockGetTargetGroupResponse:           targetgroup.TargetGroup{ID: "123"},
			mockGetTargetGroupDeploymentResponse: targetgroup.Deployment{ID: "abc"},
			want:                                 `{"createdAt":"0001-01-01T00:00:00Z","icon":"","id":"123","targetDeployments":[{"Diagnostics":null,"Id":"abc","Priority":0,"Valid":false}],"targetSchema":{"From":"","Schema":{}},"updatedAt":"0001-01-01T00:00:00Z"}`,
			deploymentId:                         "abc",
			mockCreate:                           &targetgroup.TargetGroup{ID: "123"},
			give:                                 `{"deploymentId": "abc", "priority": 100}`,
		},

		{
			name:                                 "priority cannot be out of range",
			wantCode:                             http.StatusBadRequest,
			mockGetTargetGroupResponse:           targetgroup.TargetGroup{ID: "123"},
			mockGetTargetGroupDeploymentResponse: targetgroup.Deployment{ID: "abc"},
			want:                                 `{"error":"request body has an error: doesn't match the schema: Error at \"/priority\": number must be at most 999"}`,
			deploymentId:                         "abc",
			mockCreate:                           &targetgroup.TargetGroup{ID: "123"},
			give:                                 `{"deploymentId": "abc", "priority": 1000}`,
			mockCreateErr:                        errors.New("request body has an error: doesn't match the schema: Error at \"/priority\": number must be at most 999"),
		},
		{
			name:                                 "priority cannot be under range",
			wantCode:                             http.StatusBadRequest,
			mockGetTargetGroupResponse:           targetgroup.TargetGroup{ID: "123"},
			mockGetTargetGroupDeploymentResponse: targetgroup.Deployment{ID: "abc"},
			want:                                 `{"error":"request body has an error: doesn't match the schema: Error at \"/priority\": number must be at least 0"}`,
			deploymentId:                         "abc",
			mockCreate:                           &targetgroup.TargetGroup{ID: "123"},
			give:                                 `{"deploymentId": "abc", "priority": -1}`,
			mockCreateErr:                        errors.New("request body has an error: doesn't match the schema: Error at \"/priority\": number must be at most 999"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			db := ddbmock.New(t)
			db.MockQueryWithErr(&storage.GetTargetGroup{Result: tc.mockGetTargetGroupResponse}, tc.mockGetTargetGroupErr)
			db.MockQueryWithErr(&storage.GetTargetGroupDeployment{Result: tc.mockGetTargetGroupDeploymentResponse}, tc.mockGetTargetGroupErr)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := mocks.NewMockTargetGroupService(ctrl)

			if tc.mockCreateErr == nil {
				m.EXPECT().CreateTargetGroupLink(gomock.Any(), gomock.Any(), gomock.Any()).Return(tc.mockCreate, tc.mockCreateErr)

			}

			a := API{DB: db, TargetGroupService: m}
			handler := newTestServer(t, &a)

			req, err := http.NewRequest("POST", "/api/v1/target-groups/123/link", strings.NewReader(tc.give))

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
