package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/common-fate/common-fate/pkg/api/mocks"
	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/common-fate/common-fate/pkg/service/internalidentitysvc"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
	"github.com/common-fate/ddb/ddbmock"
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
			wantBody: `{"groups":[{"description":"","id":"123","memberCount":0,"members":[],"name":"test","source":"test"},{"description":"","id":"1234","memberCount":2,"members":["1","2"],"name":"test","source":"test"}],"next":""}`,
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

func TestCreateGroup(t *testing.T) {
	type testcase struct {
		name                   string
		body                   string
		wantCode               int
		wantBody               string
		withUser               string
		expectCreateGroupOpts  *types.CreateGroupRequest
		withCreatedGroup       *identity.Group
		expectCreateGroupError error
	}

	adminGroup := "test_admins"
	testcases := []testcase{
		{name: "create internal group ok",
			body:                  `{"name":"test","description":"user","members": []}`,
			wantCode:              http.StatusCreated,
			expectCreateGroupOpts: &types.CreateGroupRequest{Name: "test", Description: aws.String("user"), Members: []string{}},
			withCreatedGroup: &identity.Group{
				ID:          "1234",
				IdpID:       "1234",
				Name:        "test",
				Description: "user",
				Users:       []string{},
				Status:      types.ACTIVE,
				Source:      "internal",
			},
			wantBody: `{"description":"user","id":"1234","memberCount":0,"members":[],"name":"test","source":"internal"}`,
		},
		{name: "create internal group user no exist error",
			body:     `{"name":"test","description":"user","members": ["123"]}`,
			wantCode: http.StatusBadRequest,

			expectCreateGroupOpts: &types.CreateGroupRequest{Name: "test", Description: aws.String("user"), Members: []string{"123"}},
			withCreatedGroup: &identity.Group{
				ID:          "1234",
				IdpID:       "1234",
				Name:        "test",
				Description: "user",
				Users:       []string{},
				Status:      types.ACTIVE,
				Source:      "internal",
			},
			wantBody:               `{"error":"user  does not exist"}`,
			expectCreateGroupError: internalidentitysvc.UserNotFoundError{},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			db := ddbmock.New(t)
			db.MockQuery(&storage.GetUser{ID: tc.withUser, Result: &identity.User{}})
			ctrl := gomock.NewController(t)

			mockIdentity := mocks.NewMockInternalIdentityService(ctrl)
			mockIdentity.EXPECT().CreateGroup(gomock.Any(), gomock.Eq(*tc.expectCreateGroupOpts)).Return(tc.withCreatedGroup, tc.expectCreateGroupError)

			a := API{AdminGroup: adminGroup, DB: db, InternalIdentity: mockIdentity}
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

func TestUpdateGroup(t *testing.T) {
	type testcase struct {
		name                   string
		body                   string
		wantCode               int
		wantBody               string
		withUser               string
		expectCreateGroupOpts  *types.CreateGroupRequest
		withExistingGroup      identity.Group
		withUpdatedGroup       *identity.Group
		existingGroupId        string
		expectCreateGroupError error
	}

	adminGroup := "test_admins"
	testcases := []testcase{
		{name: "update existing group ok",
			body:     `{"name":"updated name","description":"user","members": []}`,
			wantCode: http.StatusOK,

			expectCreateGroupOpts: &types.CreateGroupRequest{Name: "updated name", Description: aws.String("user"), Members: []string{}},
			existingGroupId:       "1234",
			withExistingGroup: identity.Group{
				ID:          "1234",
				IdpID:       "1234",
				Name:        "test",
				Description: "user",
				Users:       []string{},
				Status:      types.ACTIVE,
				Source:      "internal",
			},
			withUpdatedGroup: &identity.Group{
				ID:          "1234",
				IdpID:       "1234",
				Name:        "updated name",
				Description: "user",
				Users:       []string{},
				Status:      types.ACTIVE,
				Source:      "internal",
			},
			wantBody: `{"description":"user","id":"1234","memberCount":0,"members":[],"name":"updated name","source":"internal"}`,
		},
		{name: "update group with no user existing",
			body:     `{"name":"updated name","description":"user","members": []}`,
			wantCode: http.StatusBadRequest,

			expectCreateGroupOpts: &types.CreateGroupRequest{Name: "updated name", Description: aws.String("user"), Members: []string{}},
			existingGroupId:       "1234",
			withExistingGroup: identity.Group{
				ID:          "1234",
				IdpID:       "1234",
				Name:        "test",
				Description: "user",
				Users:       []string{},
				Status:      types.ACTIVE,
				Source:      "internal",
			},
			withUpdatedGroup: &identity.Group{
				ID:          "1234",
				IdpID:       "1234",
				Name:        "updated name",
				Description: "user",
				Users:       []string{},
				Status:      types.ACTIVE,
				Source:      "internal",
			},
			wantBody:               `{"error":"user  does not exist"}`,
			expectCreateGroupError: internalidentitysvc.UserNotFoundError{}},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			db := ddbmock.New(t)
			db.MockQuery(&storage.GetUser{ID: tc.withUser, Result: &identity.User{Groups: []string{}}})
			db.MockQuery(&storage.GetGroup{ID: tc.existingGroupId, Result: &tc.withExistingGroup})
			ctrl := gomock.NewController(t)

			mockIdentity := mocks.NewMockInternalIdentityService(ctrl)
			mockIdentity.EXPECT().UpdateGroup(gomock.Any(), tc.withExistingGroup, *tc.expectCreateGroupOpts).Return(tc.withUpdatedGroup, tc.expectCreateGroupError)

			a := API{AdminGroup: adminGroup, DB: db, InternalIdentity: mockIdentity}
			handler := newTestServer(t, &a)

			req, err := http.NewRequest("PUT", "/api/v1/admin/groups/"+tc.existingGroupId, strings.NewReader(tc.body))
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

func TestDeleteGroup(t *testing.T) {
	type testcase struct {
		name            string
		id              string
		wantCode        int
		wantBody        string
		withGroup       identity.Group
		withGroupError  error
		withDeleteError error
	}

	testcases := []testcase{
		{
			name:     "ok",
			wantCode: http.StatusOK,
			id:       "groupid",
			wantBody: `null`,
		},
		{
			name:           "group not found",
			wantCode:       http.StatusNotFound,
			id:             "groupid",
			withGroupError: ddb.ErrNoItems,
			wantBody:       `{"error":"group not found"}`,
		},
		{
			name:            "group not found",
			wantCode:        http.StatusBadRequest,
			id:              "groupid",
			withDeleteError: internalidentitysvc.ErrNotInternal,
			wantBody:        `{"error":"cannot update group because it is not an internal group"}`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			db := ddbmock.New(t)
			db.MockQueryWithErr(&storage.GetGroup{Result: &tc.withGroup}, tc.withGroupError)

			ctrl := gomock.NewController(t)

			mockIdentity := mocks.NewMockInternalIdentityService(ctrl)
			mockIdentity.EXPECT().DeleteGroup(gomock.Any(), tc.withGroup).AnyTimes().Return(tc.withDeleteError)

			a := API{DB: db, InternalIdentity: mockIdentity}
			handler := newTestServer(t, &a)

			req, err := http.NewRequest("DELETE", "/api/v1/admin/groups/"+tc.id, nil)
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
