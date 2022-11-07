package api

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/common-fate/ddb"
	"github.com/common-fate/ddb/ddbmock"
	"github.com/common-fate/granted-approvals/pkg/api/mocks"
	"github.com/common-fate/granted-approvals/pkg/identity"
	"github.com/common-fate/granted-approvals/pkg/service/cognitosvc"
	"github.com/common-fate/granted-approvals/pkg/storage"
	"github.com/common-fate/granted-approvals/pkg/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetUsers(t *testing.T) {
	type testcase struct {
		name     string
		idpErr   error
		idpUsers []identity.User
		wantCode int
		wantBody string
	}

	testcases := []testcase{
		{
			name:     "ok",
			wantCode: http.StatusOK,
			idpUsers: []identity.User{
				{
					ID:        "123",
					Email:     "test@acme.com",
					FirstName: "Test",
					LastName:  "User",
				},
				{
					ID:        "1234",
					Email:     "test2@acme.com",
					FirstName: "Test",
					LastName:  "User",
				},
			},
			wantBody: `{"next":null,"users":[{"email":"test@acme.com","firstName":"Test","groups":[],"id":"123","lastName":"User","picture":"","status":"","updatedAt":"0001-01-01T00:00:00Z"},{"email":"test2@acme.com","firstName":"Test","groups":[],"id":"1234","lastName":"User","picture":"","status":"","updatedAt":"0001-01-01T00:00:00Z"}]}`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			db := ddbmock.New(t)
			db.MockQueryWithErr(&storage.ListUsersForStatus{Result: tc.idpUsers}, tc.idpErr)

			a := API{DB: db}
			handler := newTestServer(t, &a)

			req, err := http.NewRequest("GET", "/api/v1/admin/users", nil)
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

func TestGetUser(t *testing.T) {
	type testcase struct {
		name     string
		idpErr   error
		idpUser  *identity.User
		wantCode int
		wantBody string
	}

	testcases := []testcase{
		{
			name:     "ok",
			wantCode: http.StatusOK,
			idpUser: &identity.User{
				Email:     "test@acme.com",
				ID:        "123",
				FirstName: "Test",
				LastName:  "User",
			},
			wantBody: `{"email":"test@acme.com","firstName":"Test","groups":[],"id":"123","lastName":"User","picture":"","status":"","updatedAt":"0001-01-01T00:00:00Z"}`,
		},
		{
			name:     "user not found",
			wantCode: http.StatusNotFound,
			idpErr:   ddb.ErrNoItems,

			wantBody: `{"error":"item query returned no items"}`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			db := ddbmock.New(t)
			db.MockQueryWithErr(&storage.GetUser{Result: tc.idpUser}, tc.idpErr)

			a := API{DB: db}
			handler := newTestServer(t, &a)

			req, err := http.NewRequest("GET", "/api/v1/users/123", nil)
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

func TestPostApiV1AdminUsers(t *testing.T) {
	type testcase struct {
		name                  string
		body                  string
		wantCode              int
		wantBody              string
		notEnabled            bool
		expectCreateUserOpts  *cognitosvc.CreateUserOpts
		withCreatedUser       *identity.User
		expectCreateUserError error
	}

	adminGroup := "test_admins"
	testcases := []testcase{
		{name: "Not enabled", body: `{"firstName":"test","lastName":"user","email":"chris@commonfate.io","isAdmin":false}`, wantCode: http.StatusBadRequest, notEnabled: true, wantBody: `{"error":"api not available"}`},
		{name: "ok",
			body:                 `{"firstName":"test","lastName":"user","email":"test@test.com","isAdmin":false}`,
			wantCode:             http.StatusCreated,
			expectCreateUserOpts: &cognitosvc.CreateUserOpts{FirstName: "test", LastName: "user", Email: "test@test.com", IsAdmin: false},
			withCreatedUser: &identity.User{
				ID:        "1234",
				FirstName: "test",
				LastName:  "user",
				Email:     "test@test.com",
				Groups:    []string{},
				Status:    types.IdpStatusACTIVE,
			},
			wantBody: `{"email":"test@test.com","firstName":"test","groups":[],"id":"1234","lastName":"user","picture":"","status":"ACTIVE","updatedAt":"0001-01-01T00:00:00Z"}`,
		},
		{name: "error from create user",
			body:                  `{"firstName":"test","lastName":"user","email":"test@test.com","isAdmin":true}`,
			wantCode:              http.StatusInternalServerError,
			expectCreateUserOpts:  &cognitosvc.CreateUserOpts{FirstName: "test", LastName: "user", Email: "test@test.com", IsAdmin: true},
			expectCreateUserError: errors.New("random error"),
			wantBody:              `{"error":"Internal Server Error"}`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			a := API{AdminGroup: adminGroup}
			if !tc.notEnabled {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				m := mocks.NewMockCognitoService(ctrl)
				if tc.expectCreateUserOpts != nil {
					a.Cognito = m
					m.EXPECT().CreateUser(gomock.Any(), gomock.Eq(*tc.expectCreateUserOpts)).Times(1).Return(tc.withCreatedUser, tc.expectCreateUserError)
				}
			}
			handler := newTestServer(t, &a)

			req, err := http.NewRequest("POST", "/api/v1/admin/users", strings.NewReader(tc.body))
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
