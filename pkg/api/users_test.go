package api

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/common-fate/ddb/ddbmock"
	"github.com/common-fate/granted-approvals/pkg/identity"
	"github.com/common-fate/granted-approvals/pkg/storage"
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
			wantBody: `{"next":null,"users":[{"email":"test@acme.com","firstName":"Test","id":"123","lastName":"User","picture":"","status":"","updatedAt":"0001-01-01T00:00:00Z"},{"email":"test2@acme.com","firstName":"Test","id":"1234","lastName":"User","picture":"","status":"","updatedAt":"0001-01-01T00:00:00Z"}]}`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
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

			data, err := ioutil.ReadAll(rr.Body)
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
			wantBody: `{"id":"123","firstName":"Test","lastName":"User","email":"test@acme.com","groups":null,"status":"","createdAt":"0001-01-01T00:00:00Z","updatedAt":"0001-01-01T00:00:00Z"}`,
		},
		{
			name:     "user not found",
			wantCode: http.StatusNotFound,
			idpErr:   identity.UserNotFoundError{User: "123"},
			wantBody: `{"error":"user 123 not found"}`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
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

			data, err := ioutil.ReadAll(rr.Body)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tc.wantBody, string(data))
		})
	}
}
