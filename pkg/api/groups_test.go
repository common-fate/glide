package api

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/common-fate/ddb/ddbmock"
	"github.com/common-fate/granted-approvals/pkg/api/mocks"
	"github.com/common-fate/granted-approvals/pkg/identity"
	"github.com/common-fate/granted-approvals/pkg/service/cognitosvc"
	"github.com/common-fate/granted-approvals/pkg/storage"
	"github.com/common-fate/granted-approvals/pkg/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestListGroups(t *testing.T) {
	type testcase struct {
		name      string
		idpGroups []identity.Group
		wantCode  int
		wantBody  string
	}

	testcases := []testcase{
		{
			name:     "ok",
			wantCode: http.StatusOK,
			idpGroups: []identity.Group{
				{
					ID:    "123",
					Name:  "test",
					Users: nil,
				},
				{
					ID:    "1234",
					Name:  "test",
					Users: []string{"1", "2"},
				},
			},
			wantBody: `{"groups":[{"description":"","id":"123","memberCount":0,"name":"test"},{"description":"","id":"1234","memberCount":2,"name":"test"}],"next":null}`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			db := ddbmock.New(t)
			db.MockQuery(&storage.ListGroupsForStatus{Result: tc.idpGroups})

			a := API{DB: db}
			handler := newTestServer(t, &a)

			req, err := http.NewRequest("GET", "/api/v1/admin/groups", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("Cognito", "approvals:admin")

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

func TestGetGroup(t *testing.T) {
	type testcase struct {
		name     string
		idpErr   error
		idpGroup *identity.Group
		wantCode int
		wantBody string
	}

	testcases := []testcase{
		{
			name:     "ok",
			wantCode: http.StatusOK,
			idpGroup: &identity.Group{

				ID:          "123",
				Name:        "Test",
				Description: "hello",
				Users:       []string{"one", "two", "three"},
			},
			wantBody: `{"description":"hello","id":"123","memberCount":3,"name":"Test"}`,
		},
		{
			name:     "group not found",
			wantCode: http.StatusNotFound,
			idpErr:   identity.UserNotFoundError{User: "123"},
			wantBody: `{"error":"user 123 not found"}`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			db := ddbmock.New(t)
			db.MockQueryWithErr(&storage.GetGroup{Result: tc.idpGroup}, tc.idpErr)

			a := API{DB: db}
			handler := newTestServer(t, &a)

			req, err := http.NewRequest("GET", "/api/v1/admin/groups/123", nil)
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

func TestPostApiV1AdminGroups(t *testing.T) {
	type testcase struct {
		name                   string
		body                   string
		wantCode               int
		wantBody               string
		notEnabled             bool
		expectCreateGroupOpts  *cognitosvc.CreateGroupOpts
		withCreatedGroup       *identity.Group
		expectCreateGroupError error
	}

	adminGroup := "test_admins"
	testcases := []testcase{
		{name: "Not enabled", body: `{"name":"test","description":"user"}`, wantCode: http.StatusBadRequest, notEnabled: true, wantBody: `{"error":"api not available"}`},
		{name: "ok",
			body:                  `{"name":"test","description":"user"}`,
			wantCode:              http.StatusCreated,
			expectCreateGroupOpts: &cognitosvc.CreateGroupOpts{Name: "test", Description: "user"},
			withCreatedGroup: &identity.Group{
				ID:          "1234",
				IdpID:       "1234",
				Name:        "test",
				Description: "user",
				Users:       []string{},
				Status:      types.IdpStatusACTIVE,
			},
			wantBody: `{"description":"user","id":"1234","memberCount":0,"name":"test"}`,
		},
		{name: "error from create user",
			body:                   `{"name":"test","description":"user"}`,
			wantCode:               http.StatusInternalServerError,
			expectCreateGroupOpts:  &cognitosvc.CreateGroupOpts{Name: "test", Description: "user"},
			expectCreateGroupError: errors.New("random error"),
			wantBody:               `{"error":"Internal Server Error"}`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {

			a := API{AdminGroup: adminGroup}
			if !tc.notEnabled {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				m := mocks.NewMockCognitoService(ctrl)
				if tc.expectCreateGroupOpts != nil {
					a.Cognito = m
					m.EXPECT().CreateGroup(gomock.Any(), gomock.Eq(*tc.expectCreateGroupOpts)).Times(1).Return(tc.withCreatedGroup, tc.expectCreateGroupError)
				}
			}
			handler := newTestServer(t, &a)

			req, err := http.NewRequest("POST", "/api/v1/admin/groups", strings.NewReader(tc.body))
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
