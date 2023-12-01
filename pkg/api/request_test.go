package api

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/api/mocks"
	"github.com/common-fate/common-fate/pkg/cache"
	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/common-fate/common-fate/pkg/service/accesssvc"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
	"github.com/common-fate/ddb/ddbmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUserCreateRequest(t *testing.T) {
	type testcase struct {
		name string
		give string
		// mockCreate    []accesssvc.CreateRequestResult
		mockCreateErr error
		wantCode      int
		wantBody      string
	}

	testcases := []testcase{
		{
			name:     "ok",
			give:     `{"timing":{"durationSeconds": 10}, "accessRuleId": "rul_123"}`,
			wantCode: http.StatusOK,
			wantBody: `null`,
		},
		{
			name:     "no duration",
			give:     `{"accessRuleId": "rul_123"}`,
			wantCode: http.StatusBadRequest,
			wantBody: `{"error":"request body has an error: doesn't match the schema: Error at \"/timing\": property \"timing\" is missing"}`,
		},
		{
			name:     "no rule ID",
			give:     `{"timing":{"durationSeconds": 10}}`,
			wantCode: http.StatusBadRequest,
			wantBody: `{"error":"request body has an error: doesn't match the schema: Error at \"/accessRuleId\": property \"accessRuleId\" is missing"}`,
		},
		{
			name:          "rule not found",
			give:          `{"timing":{"durationSeconds": 10}, "accessRuleId": "rul_123"}`,
			mockCreateErr: apio.NewRequestError(accesssvc.ErrRuleNotFound, http.StatusBadRequest),
			wantCode:      http.StatusBadRequest,
			wantBody:      `{"error":"access rule not found"}`,
		},
		{
			name:          "no matching group",
			give:          `{"timing":{"durationSeconds": 10}, "accessRuleId": "rul_123"}`,
			mockCreateErr: apio.NewRequestError(accesssvc.ErrNoMatchingGroup, http.StatusBadRequest),
			wantCode:      http.StatusBadRequest,
			wantBody:      `{"error":"user was not in a matching group for the access rule"}`,
		},
	}

	for i := range testcases {
		tc := testcases[i]
		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			mockAccess := mocks.NewMockAccessService(ctrl)
			mockAccess.EXPECT().CreateRequests(gomock.Any(), gomock.Any()).Return(nil, tc.mockCreateErr).AnyTimes()
			a := API{Access: mockAccess}
			handler := newTestServer(t, &a)

			req, err := http.NewRequest("POST", "/api/v1/requests", strings.NewReader(tc.give))
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

func TestUserCancelRequest(t *testing.T) {
	type testcase struct {
		name          string
		give          string
		mockCancelErr error
		wantCode      int
		wantBody      string
	}

	testcases := []testcase{
		{
			name:          "ok",
			give:          `{}`,
			mockCancelErr: nil,
			wantCode:      http.StatusOK,
			wantBody:      `{}`,
		},
		{
			name:          "unauthorized",
			give:          `{}`,
			mockCancelErr: accesssvc.ErrUserNotAuthorized,
			wantCode:      http.StatusUnauthorized,
			wantBody:      `{"error":"user is not authorized to perform this action"}`,
		},
		{
			name:          "not found",
			give:          `{}`,
			mockCancelErr: ddb.ErrNoItems,
			wantCode:      http.StatusNotFound,
			wantBody:      `{"error":"item query returned no items"}`,
		},
		{
			name:          "cannot be cancelled",
			give:          `{}`,
			mockCancelErr: accesssvc.ErrRequestCannotBeCancelled,
			wantCode:      http.StatusBadRequest,
			wantBody:      `{"error":"only pending requests can be cancelled"}`,
		},
		{
			name:          "unhandled dynamo db error",
			give:          `{}`,
			mockCancelErr: errors.New("an error we don't handle"),
			wantCode:      http.StatusInternalServerError,
			wantBody:      `{"error":"Internal Server Error"}`,
		},
	}

	for i := range testcases {
		tc := testcases[i]
		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			mockAccess := mocks.NewMockAccessService(ctrl)
			mockAccess.EXPECT().CancelRequest(gomock.Any(), gomock.Any()).Return(tc.mockCancelErr).AnyTimes()
			a := API{Access: mockAccess}
			handler := newTestServer(t, &a)

			req, err := http.NewRequest("POST", "/api/v1/requests/abcd/cancel", strings.NewReader(tc.give))
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

func TestUserGetRequest(t *testing.T) {

	type testcase struct {
		name string
		// A stringified JSON object that will be used as the request body (should be empty for these GET Requests)
		givenID                  string
		mockGetRequestErr        error
		mockGetReviewerErr       error
		mockGetAccessRuleVersion *rule.AccessRule
		// request body (Request type)
		mockGetRequest *access.Request
		// request body (Request type)
		mockGetReviewer *access.Reviewer
		// expected HTTP response code
		wantCode int
		// expected HTTP response body
		wantBody                     string
		withRequestArgumentsResponse map[string]types.RequestArgument
	}

	testcases := []testcase{
		{
			name:     "requestor can see their own request",
			givenID:  `req_123`,
			wantCode: http.StatusOK,
			mockGetRequest: &access.Request{
				ID:          "req_123",
				Status:      access.PENDING,
				Rule:        "abcd",
				RuleVersion: "efgh",
			},
			mockGetAccessRuleVersion:     &rule.AccessRule{ID: "test"},
			withRequestArgumentsResponse: make(map[string]types.RequestArgument),
			// canReview is false in the response
			wantBody: `{"accessRule":{"createdAt":"0001-01-01T00:00:00Z","description":"","id":"test","isCurrent":false,"name":"","target":{"provider":{"id":"","type":""}},"timeConstraints":{"maxDurationSeconds":0},"updatedAt":"0001-01-01T00:00:00Z","version":""},"arguments":{},"canReview":false,"id":"req_123","requestedAt":"0001-01-01T00:00:00Z","requestor":"","status":"PENDING","timing":{"durationSeconds":0},"updatedAt":"0001-01-01T00:00:00Z"}`,
		},
		{
			name:     "reviewer can see request they can review",
			givenID:  `req_123`,
			wantCode: http.StatusOK,
			mockGetRequest: &access.Request{
				RequestedBy: "randomUser",
				ID:          "req_123",
				Status:      access.PENDING,
				Rule:        "abcd",
				RuleVersion: "efgh",
			},
			mockGetAccessRuleVersion:     &rule.AccessRule{ID: "test"},
			withRequestArgumentsResponse: make(map[string]types.RequestArgument),
			mockGetReviewer: &access.Reviewer{Request: access.Request{
				ID:          "req_123",
				Status:      access.PENDING,
				Rule:        "abcd",
				RuleVersion: "efgh",
			}},
			// note canReview is true in the response
			wantBody: `{"accessRule":{"createdAt":"0001-01-01T00:00:00Z","description":"","id":"test","isCurrent":false,"name":"","target":{"provider":{"id":"","type":""}},"timeConstraints":{"maxDurationSeconds":0},"updatedAt":"0001-01-01T00:00:00Z","version":""},"arguments":{},"canReview":true,"id":"req_123","requestedAt":"0001-01-01T00:00:00Z","requestor":"","status":"PENDING","timing":{"durationSeconds":0},"updatedAt":"0001-01-01T00:00:00Z"}`,
		},
		{
			name:              "noRequestFound",
			givenID:           `wrongID`,
			wantCode:          http.StatusNotFound,
			mockGetRequestErr: ddb.ErrNoItems,
			wantBody:          `{"error":"item query returned no items"}`,
		},
		{
			name:    "not a requestor or reviewer",
			givenID: `req_123`,
			mockGetRequest: &access.Request{
				ID:          "req_123",
				RequestedBy: "notThisUser",
				Status:      access.PENDING,
			},
			withRequestArgumentsResponse: make(map[string]types.RequestArgument),
			mockGetAccessRuleVersion:     &rule.AccessRule{ID: "test"},
			mockGetReviewerErr:           ddb.ErrNoItems,
			wantCode:                     http.StatusNotFound,
			wantBody:                     `{"error":"item query returned no items"}`,
		},
	}

	for i := range testcases {
		tc := testcases[i]
		t.Run(tc.name, func(t *testing.T) {

			db := ddbmock.New(t)
			db.MockQueryWithErr(&storage.GetRequest{Result: tc.mockGetRequest}, tc.mockGetRequestErr)
			db.MockQueryWithErr(&storage.GetRequestReviewer{Result: tc.mockGetReviewer}, tc.mockGetReviewerErr)
			db.MockQuery(&storage.GetAccessRuleVersion{Result: tc.mockGetAccessRuleVersion})
			db.MockQuery(&storage.ListCachedProviderOptions{Result: []cache.ProviderOption{}})
			ctrl := gomock.NewController(t)
			rs := mocks.NewMockAccessRuleService(ctrl)
			if tc.withRequestArgumentsResponse != nil {
				rs.EXPECT().RequestArguments(gomock.Any(), gomock.Any()).Return(tc.withRequestArgumentsResponse, nil)
			}
			a := API{DB: db, Rules: rs}
			handler := newTestServer(t, &a)

			req, err := http.NewRequest("GET", "/api/v1/requests/"+tc.givenID, strings.NewReader(""))
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
			} else {
				fmt.Print((data))
			}

			if tc.wantBody != "" {
				assert.Equal(t, tc.wantBody, string(data))
			}
		})
	}

}

func TestUserListRequests(t *testing.T) {

	type testcase struct {
		name           string
		giveReviewer   *string
		giveStatus     *string
		mockDBQuery    ddb.QueryBuilder
		mockDBQueryErr error
		// expected HTTP response code
		wantCode int
		// expected HTTP response body
		wantBody string
	}
	approved := "APPROVED"
	reviewer := "true"
	badReviewer := "hello"
	badStatus := "hi"
	testcases := []testcase{

		{
			name:     "ok requestor",
			wantCode: http.StatusOK,
			mockDBQuery: &storage.ListRequestsForUser{Result: []access.Request{{
				ID:          "req_123",
				Status:      access.PENDING,
				Rule:        "abcd",
				RuleVersion: "efgh",
			}}},

			wantBody: `{"next":null,"requests":[{"accessRuleId":"abcd","accessRuleVersion":"efgh","id":"req_123","requestedAt":"0001-01-01T00:00:00Z","requestor":"","status":"PENDING","timing":{"durationSeconds":0},"updatedAt":"0001-01-01T00:00:00Z"}]}`,
		},
		{
			name:       "ok requestor with status",
			wantCode:   http.StatusOK,
			giveStatus: &approved,
			mockDBQuery: &storage.ListRequestsForUserAndStatus{Result: []access.Request{{
				ID:          "req_123",
				Status:      access.APPROVED,
				Rule:        "abcd",
				RuleVersion: "efgh",
			}}},

			wantBody: `{"next":null,"requests":[{"accessRuleId":"abcd","accessRuleVersion":"efgh","id":"req_123","requestedAt":"0001-01-01T00:00:00Z","requestor":"","status":"APPROVED","timing":{"durationSeconds":0},"updatedAt":"0001-01-01T00:00:00Z"}]}`,
		},
		{
			name:         "ok reviewer",
			wantCode:     http.StatusOK,
			giveReviewer: &reviewer,
			mockDBQuery: &storage.ListRequestsForReviewer{Result: []access.Request{{
				ID:          "req_123",
				Status:      access.PENDING,
				Rule:        "abcd",
				RuleVersion: "efgh",
			}}},

			wantBody: `{"next":null,"requests":[{"accessRuleId":"abcd","accessRuleVersion":"efgh","id":"req_123","requestedAt":"0001-01-01T00:00:00Z","requestor":"","status":"PENDING","timing":{"durationSeconds":0},"updatedAt":"0001-01-01T00:00:00Z"}]}`,
		},
		{
			name:         "ok requestor with status",
			wantCode:     http.StatusOK,
			giveStatus:   &approved,
			giveReviewer: &reviewer,
			mockDBQuery: &storage.ListRequestsForReviewerAndStatus{Result: []access.Request{{
				ID:          "req_123",
				Status:      access.APPROVED,
				Rule:        "abcd",
				RuleVersion: "efgh",
			}}},

			wantBody: `{"next":null,"requests":[{"accessRuleId":"abcd","accessRuleVersion":"efgh","id":"req_123","requestedAt":"0001-01-01T00:00:00Z","requestor":"","status":"APPROVED","timing":{"durationSeconds":0},"updatedAt":"0001-01-01T00:00:00Z"}]}`,
		},
		{
			name:           "internal error user",
			wantCode:       http.StatusInternalServerError,
			mockDBQuery:    &storage.ListRequestsForUser{},
			mockDBQueryErr: errors.New("random error"),
			wantBody:       `{"error":"Internal Server Error"}`,
		},
		{
			name:           "internal error user and status",
			wantCode:       http.StatusInternalServerError,
			giveStatus:     &approved,
			mockDBQuery:    &storage.ListRequestsForUserAndStatus{},
			mockDBQueryErr: errors.New("random error"),
			wantBody:       `{"error":"Internal Server Error"}`,
		},
		{
			name:           "internal error reviewer",
			wantCode:       http.StatusInternalServerError,
			giveReviewer:   &reviewer,
			mockDBQuery:    &storage.ListRequestsForReviewer{},
			mockDBQueryErr: errors.New("random error"),
			wantBody:       `{"error":"Internal Server Error"}`,
		},
		{
			name:           "internal error reviewer and status",
			wantCode:       http.StatusInternalServerError,
			giveReviewer:   &reviewer,
			giveStatus:     &approved,
			mockDBQuery:    &storage.ListRequestsForReviewerAndStatus{},
			mockDBQueryErr: errors.New("random error"),
			wantBody:       `{"error":"Internal Server Error"}`,
		},
		{
			name:         "bad reviewer param",
			wantCode:     http.StatusBadRequest,
			giveReviewer: &badReviewer,

			wantBody: `{"error":"parameter \"reviewer\" in query has an error: value hello: an invalid boolean: invalid syntax"}`,
		},
		{
			name:         "bad status",
			wantCode:     http.StatusBadRequest,
			giveReviewer: &reviewer,
			giveStatus:   &badStatus,
			wantBody:     `{"error":"parameter \"status\" in query has an error: value is not one of the allowed values"}`,
		},
	}

	for i := range testcases {
		tc := testcases[i]
		t.Run(tc.name, func(t *testing.T) {

			db := ddbmock.New(t)
			db.MockQueryWithErr(tc.mockDBQuery, tc.mockDBQueryErr)
			a := API{DB: db}
			handler := newTestServer(t, &a)
			var qp []string
			if tc.giveReviewer != nil {
				qp = append(qp, "reviewer="+*tc.giveReviewer)
			}
			if tc.giveStatus != nil {
				qp = append(qp, "status="+*tc.giveStatus)
			}

			req, err := http.NewRequest("GET", "/api/v1/requests?"+strings.Join(qp, "&"), strings.NewReader(""))
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
			} else {
				fmt.Print((data))
			}

			if tc.wantBody != "" {
				assert.Equal(t, tc.wantBody, string(data))
			}
		})
	}

}
func TestRevokeRequest(t *testing.T) {
	type testcase struct {
		request              access.Request
		name                 string
		give                 string
		withUID              string
		withUEmail           string
		withRevokeGrantErr   error
		wantCode             int
		withGetRequestError  error
		withGetReviewerError error
		wantBody             string
		withIsAdmin          bool
	}

	testcases := []testcase{

		{
			name:                "grant not found",
			request:             access.Request{},
			wantCode:            http.StatusNotFound,
			withGetRequestError: ddb.ErrNoItems,
			wantBody:            `{"error":"request not found or you don't have access to it"}`,
		},
		{
			name: "user can revoke their own grant",
			request: access.Request{
				RequestedBy: "user1",
			},
			withUID:    "user1",
			withUEmail: "user1@mail.com",
			wantCode:   http.StatusOK,
		},
		{
			name: "admin can revoke any request",
			request: access.Request{
				RequestedBy: "user1",
			},
			withUID:     "admin",
			withUEmail:  "admin@mail.com",
			withIsAdmin: true,
			wantCode:    http.StatusOK,
		},
		{
			name: "user cant revoke other users request",
			request: access.Request{
				RequestedBy: "user1",
			},
			withUID:              "user2",
			withUEmail:           "user2@mail.com",
			withGetReviewerError: ddb.ErrNoItems,
			wantCode:             http.StatusNotFound,
			wantBody:             `{"error":"request not found or you don't have access to it"}`,
		},
		{
			name: "reviewer can revoke request",
			request: access.Request{
				RequestedBy: "user1",
			},
			withUID:    "user2",
			withUEmail: "user2@mail.com",
			wantCode:   http.StatusOK,
		},
	}

	for i := range testcases {
		tc := testcases[i]
		t.Run(tc.name, func(t *testing.T) {

			db := ddbmock.New(t)
			db.MockQueryWithErr(&storage.GetRequest{Result: &tc.request}, tc.withGetRequestError)
			db.MockQueryWithErr(&storage.GetRequestReviewer{Result: &access.Reviewer{Request: tc.request}}, tc.withGetReviewerError)
			ctrl := gomock.NewController(t)
			workflowMock := mocks.NewMockWorkflow(ctrl)
			workflowMock.EXPECT().Revoke(gomock.Any(), tc.request, tc.withUID, tc.withUEmail).AnyTimes().Return(nil, tc.withRevokeGrantErr)

			a := API{DB: db, Workflow: workflowMock}
			handler := newTestServer(t, &a, withIsAdmin(tc.withIsAdmin), withRequestUser(identity.User{ID: tc.withUID, Email: tc.withUEmail}))

			req, err := http.NewRequest("POST", "/api/v1/requests/123/revoke", strings.NewReader(tc.give))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Add("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)
			data, err := io.ReadAll(rr.Body)
			if err != nil {
				t.Fatal(err)
			} else {
				fmt.Print((data))
			}
			assert.Equal(t, tc.wantCode, rr.Code)
			if tc.wantBody != "" {
				assert.Equal(t, tc.wantBody, string(data))
			}
		})
	}
}
func TestUserListRequestEvents(t *testing.T) {

	type testcase struct {
		name                      string
		mockGetRequest            storage.GetRequest
		mockGetRequestErr         error
		mockGetRequestReviewer    storage.GetRequestReviewer
		mockGetRequestReviewerErr error
		mockListEvents            storage.ListRequestEvents
		mockListEventsErr         error
		apiUserID                 string
		apiUserIsAdmin            bool
		// expected HTTP response code
		wantCode int
		// expected HTTP response body
		wantBody string
	}

	testcases := []testcase{

		{
			name:     "ok requestor",
			wantCode: http.StatusOK,
			mockGetRequest: storage.GetRequest{
				ID: "1234",
				Result: &access.Request{
					ID:          "1234",
					RequestedBy: "abcd",
				},
			},
			mockListEvents: storage.ListRequestEvents{
				RequestID: "1234",
				Result: []access.RequestEvent{
					{ID: "event", RequestID: "1234"},
				},
			},
			apiUserID: "abcd",
			wantBody:  `{"events":[{"createdAt":"0001-01-01T00:00:00Z","id":"event","requestId":"1234"}],"next":null}`,
		},
		{
			name:     "ok reviewer",
			wantCode: http.StatusOK,
			mockGetRequest: storage.GetRequest{
				ID: "1234",
				Result: &access.Request{
					ID:          "1234",
					RequestedBy: "wrong",
				},
			},
			mockGetRequestReviewer: storage.GetRequestReviewer{
				RequestID:  "1234",
				ReviewerID: "abcd",
				Result: &access.Reviewer{
					ReviewerID: "abcd",
					Request: access.Request{
						ID:          "1234",
						RequestedBy: "wrong",
					},
				},
			},
			mockListEvents: storage.ListRequestEvents{
				RequestID: "1234",
				Result: []access.RequestEvent{
					{ID: "event", RequestID: "1234"},
				},
			},
			apiUserID: "abcd",
			wantBody:  `{"events":[{"createdAt":"0001-01-01T00:00:00Z","id":"event","requestId":"1234"}],"next":null}`,
		},
		{
			name:     "ok admin",
			wantCode: http.StatusOK,
			mockGetRequest: storage.GetRequest{
				ID: "1234",
				Result: &access.Request{
					ID:          "1234",
					RequestedBy: "wrong",
				},
			},
			mockGetRequestReviewerErr: ddb.ErrNoItems,
			mockListEvents: storage.ListRequestEvents{
				RequestID: "1234",
				Result: []access.RequestEvent{
					{ID: "event", RequestID: "1234"},
				},
			},
			apiUserID:      "abcd",
			apiUserIsAdmin: true,
			wantBody:       `{"events":[{"createdAt":"0001-01-01T00:00:00Z","id":"event","requestId":"1234"}],"next":null}`,
		},
		{
			name:              "not found",
			wantCode:          http.StatusUnauthorized,
			mockGetRequestErr: ddb.ErrNoItems,

			wantBody: `{"error":"item query returned no items"}`,
		},
	}

	for i := range testcases {
		tc := testcases[i]
		t.Run(tc.name, func(t *testing.T) {

			db := ddbmock.New(t)
			db.MockQueryWithErr(&tc.mockGetRequest, tc.mockGetRequestErr)
			db.MockQueryWithErr(&tc.mockListEvents, tc.mockListEventsErr)
			db.MockQueryWithErr(&tc.mockGetRequestReviewer, tc.mockGetRequestReviewerErr)
			a := API{DB: db}
			handler := newTestServer(t, &a, withRequestUser(identity.User{ID: tc.apiUserID}), withIsAdmin(tc.apiUserIsAdmin))

			req, err := http.NewRequest("GET", "/api/v1/requests/1234/events", strings.NewReader(""))
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
			} else {
				fmt.Print((data))
			}

			if tc.wantBody != "" {
				assert.Equal(t, tc.wantBody, string(data))
			}
		})
	}

}

func TestGetAccessToken(t *testing.T) {
	type testcase struct {
		name               string
		withUID            string
		withRequest        *access.Request
		withGetRequestErr  error
		withAccessToken    *access.AccessToken
		withAccessTokenErr error
		wantBody           string
		wantCode           int
	}

	testcases := []testcase{
		{
			name:            "ok",
			withUID:         "a",
			withRequest:     &access.Request{RequestedBy: "a"},
			withAccessToken: &access.AccessToken{Token: "token"},
			wantBody:        `{"hasToken":true,"token":"token"}`,
			wantCode:        http.StatusOK,
		},
		{
			name:            "wrong user unauthorised",
			withUID:         "b",
			withRequest:     &access.Request{RequestedBy: "a"},
			withAccessToken: &access.AccessToken{Token: "token"},
			wantBody:        `{"error":"not authorised"}`,
			wantCode:        http.StatusUnauthorized,
		},
		{
			name:              "request not found",
			withUID:           "b",
			withGetRequestErr: ddb.ErrNoItems,
			withAccessToken:   &access.AccessToken{Token: "token"},
			wantBody:          `{"error":"request not found"}`,
			wantCode:          http.StatusNotFound,
		},
		{
			name:               "request has no token",
			withUID:            "a",
			withRequest:        &access.Request{RequestedBy: "a"},
			withAccessTokenErr: ddb.ErrNoItems,
			wantBody:           `{"hasToken":false}`,
			wantCode:           http.StatusOK,
		},
	}

	for i := range testcases {
		tc := testcases[i]
		t.Run(tc.name, func(t *testing.T) {

			db := ddbmock.New(t)
			db.MockQueryWithErr(&storage.GetRequest{Result: tc.withRequest}, tc.withGetRequestErr)
			db.MockQueryWithErr(&storage.GetAccessToken{Result: tc.withAccessToken}, tc.withAccessTokenErr)

			a := API{DB: db}
			handler := newTestServer(t, &a, withRequestUser(identity.User{ID: tc.withUID}))

			req, err := http.NewRequest("GET", "/api/v1/requests/123/access-token", strings.NewReader(""))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Add("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)
			data, err := io.ReadAll(rr.Body)
			if err != nil {
				t.Fatal(err)
			} else {
				fmt.Print((data))
			}
			assert.Equal(t, tc.wantCode, rr.Code)
			if tc.wantBody != "" {
				assert.Equal(t, tc.wantBody, string(data))
			}
		})
	}
}
