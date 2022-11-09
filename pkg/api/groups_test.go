package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/common-fate/ddb"
	"github.com/common-fate/ddb/ddbmock"
	"github.com/common-fate/granted-approvals/pkg/identity"
	"github.com/common-fate/granted-approvals/pkg/storage"
	"github.com/common-fate/granted-approvals/pkg/types"
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
					ID:     "123",
					Name:   "test",
					Users:  nil,
					Source: "test",
				},
				{
					ID:     "1234",
					Name:   "test",
					Users:  []string{"1", "2"},
					Source: "test",
				},
			},
			wantBody: `{"groups":[{"description":"","id":"123","memberCount":0,"members":null,"name":"test","source":"test"},{"description":"","id":"1234","memberCount":2,"members":["1","2"],"name":"test","source":"test"}],"next":""}`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
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
			idpErr:   ddb.ErrNoItems,

			wantBody: `{"error":"item query returned no items"}`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
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
		name                  string
		body                  string
		wantCode              int
		wantBody              string
		withUser              string
		expectCreateGroupOpts *types.CreateGroupRequest
		withCreatedGroup      *identity.Group
		// expectCreateGroupError error
	}

	adminGroup := "test_admins"
	testcases := []testcase{
		{name: "ok",
			body:     `{"id": "1234", "name":"test","description":"user","members": []}`,
			wantCode: http.StatusCreated,

			expectCreateGroupOpts: &types.CreateGroupRequest{Name: "test", Description: aws.String("user"), Members: []string{}},
			withCreatedGroup: &identity.Group{
				ID:          "1234",
				IdpID:       "1234",
				Name:        "test",
				Description: "user",
				Users:       []string{},
				Status:      types.IdpStatusACTIVE,
				Source:      "test",
			},
			wantBody: `{"description":"user","id":"1234","memberCount":0,"members":[],"name":"test","source":"internal"}`,
		},
		{name: "users added to group",
			body:                  `{"id": "1234", "name":"test","description":"user","members": ["user_1"]}`,
			wantCode:              http.StatusCreated,
			expectCreateGroupOpts: &types.CreateGroupRequest{Name: "test", Description: aws.String("user"), Members: []string{"user_1"}},
			withCreatedGroup: &identity.Group{
				ID:          "1234",
				IdpID:       "1234",
				Name:        "test",
				Description: "user",
				Users:       []string{"user_1"},
				Status:      types.IdpStatusACTIVE,
				Source:      "test",
			},
			wantBody: `{"description":"user","id":"1234","memberCount":1,"members":["user_1"],"name":"test","source":"internal"}`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			db := ddbmock.New(t)
			db.MockQuery(&storage.GetUser{ID: tc.withUser, Result: &identity.User{Groups: []string{}}})
			a := API{AdminGroup: adminGroup, DB: db}

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
