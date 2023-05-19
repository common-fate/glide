package api

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/api/mocks"
	"github.com/common-fate/common-fate/pkg/cache"
	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/common-fate/common-fate/pkg/service/accesssvc"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/target"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
	"github.com/common-fate/ddb/ddbmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// big request sample to be used in all different tests
var r = access.RequestWithGroupsWithTargets{
	Request: access.Request{
		ID:               "req_123",
		RequestStatus:    "ACTIVE",
		GroupTargetCount: 4,
		RequestedBy: access.RequestedBy{
			Email:     "user1@gmail.com",
			FirstName: "",
			ID:        "",
			LastName:  "",
		},
		Purpose: access.Purpose{Reason: aws.String("sample reason")},
	},
	Groups: []access.GroupWithTargets{
		{
			Group: access.Group{
				ID:        "group1",
				RequestID: "req_123",
				AccessRuleSnapshot: rule.AccessRule{
					ID:          "234",
					Priority:    99,
					Metadata:    rule.AccessRuleMetadata{},
					Name:        "test",
					Description: " ",
					Targets: []rule.Target{
						{
							TargetGroup: target.Group{
								ID: "aws",
							},
						},
					},
					TimeConstraints: types.AccessRuleTimeConstraints{
						MaxDurationSeconds:     3600,
						DefaultDurationSeconds: 3600,
					},
					Groups:   []string{"123"},
					Approval: rule.Approval{},
				},
				Status:        types.RequestAccessGroupStatusAPPROVED,
				RequestStatus: types.ACTIVE,
				RequestedTiming: access.Timing{
					Duration: 3600,
				},

				RequestedBy: access.RequestedBy{
					Email:     "abc@gmail.com",
					FirstName: "abc",
					ID:        "123",
					LastName:  "xyz",
				},
			},
			Targets: []access.GroupTarget{
				{
					ID:            "xyz",
					GroupID:       "group1",
					RequestID:     "req_123",
					RequestStatus: "ACTIVE",
					RequestedBy: access.RequestedBy{
						Email:     "abc@gmail.com",
						FirstName: "abc",
						ID:        "123",
						LastName:  "xyz",
					},
					TargetCacheID: "",
					TargetGroupID: "aws",
					TargetKind: cache.Kind{
						Publisher: "",
						Name:      "",
						Kind:      "",
						Icon:      "",
					},
					Fields: []access.Field{
						{
							ID: "123",
						},
					},
					Grant: &access.Grant{
						Subject: "person",
						Status:  "ACTIVE",
					},
				},
			},
		},
		{
			Group: access.Group{
				ID:        "group2",
				RequestID: "123",
				AccessRuleSnapshot: rule.AccessRule{
					ID:          "234",
					Priority:    99,
					Metadata:    rule.AccessRuleMetadata{},
					Name:        "test",
					Description: " ",
					Targets: []rule.Target{
						{
							TargetGroup: target.Group{
								ID: "aws",
							},
						},
					},
					TimeConstraints: types.AccessRuleTimeConstraints{
						MaxDurationSeconds:     3600,
						DefaultDurationSeconds: 3600,
					},
					Groups:   []string{"123"},
					Approval: rule.Approval{},
				},
				Status:        types.RequestAccessGroupStatusAPPROVED,
				RequestStatus: types.ACTIVE,
				RequestedTiming: access.Timing{
					Duration: 3600,
				},

				RequestedBy: access.RequestedBy{
					Email:     "",
					FirstName: "",
					ID:        "",
					LastName:  "",
				},
			},
			Targets: []access.GroupTarget{
				{
					ID:            "xyz",
					GroupID:       "group2",
					RequestID:     "req_123",
					RequestStatus: "ACTIVE",
					RequestedBy: access.RequestedBy{
						Email:     "",
						FirstName: "",
						ID:        "",
						LastName:  "",
					},
					TargetCacheID: "",
					TargetGroupID: "aws",
					TargetKind: cache.Kind{
						Publisher: "",
						Name:      "",
						Kind:      "",
						Icon:      "",
					},
					Fields: []access.Field{
						{
							ID: "123",
						},
					},
					Grant: &access.Grant{
						Subject: "person",
						Status:  "ACTIVE",
					},
				},
			},
		},
	},
}

func TestUserCreateRequest(t *testing.T) {
	type testcase struct {
		name          string
		give          string
		mockCreateErr error
		wantCode      int
		wantBody      string

		createRes access.RequestWithGroupsWithTargets
	}

	testcases := []testcase{
		{
			name: "ok",
			give: `{
  "groupOptions": [
    {
      "id": "group1",
      "timing": {
        "durationSeconds": 3600
      }
    },
    {
      "id": "group2",
      "timing": {
        "durationSeconds": 1800
      }
    }
  ],
  "createTemplate": false,
  "preflightId": "1234567890",
  "reason": "Sample reason"
}`,
			wantCode:  http.StatusOK,
			wantBody:  `{"accessGroups":[{"accessRule":{"timeConstraints":{"maxDurationSeconds":0}},"createdAt":"0001-01-01T00:00:00Z","id":"group1","requestId":"req_123","requestStatus":"ACTIVE","requestedBy":{"email":"abc@gmail.com","firstName":"abc","id":"123","lastName":"xyz"},"requestedTiming":{"durationSeconds":0},"status":"APPROVED","targets":[{"accessGroupId":"group1","fields":[{"fieldTitle":"","id":"123","value":"","valueLabel":""}],"id":"xyz","requestId":"req_123","requestedBy":{"email":"abc@gmail.com","firstName":"abc","id":"123","lastName":"xyz"},"status":"ACTIVE","targetGroupId":"aws","targetKind":{"icon":"","kind":"","name":"","publisher":""}}],"updatedAt":"0001-01-01T00:00:00Z"},{"accessRule":{"timeConstraints":{"maxDurationSeconds":0}},"createdAt":"0001-01-01T00:00:00Z","id":"group2","requestId":"123","requestStatus":"ACTIVE","requestedBy":{"email":"","firstName":"","id":"","lastName":""},"requestedTiming":{"durationSeconds":0},"status":"APPROVED","targets":[{"accessGroupId":"group2","fields":[{"fieldTitle":"","id":"123","value":"","valueLabel":""}],"id":"xyz","requestId":"req_123","requestedBy":{"email":"","firstName":"","id":"","lastName":""},"status":"ACTIVE","targetGroupId":"aws","targetKind":{"icon":"","kind":"","name":"","publisher":""}}],"updatedAt":"0001-01-01T00:00:00Z"}],"id":"req_123","purpose":{"reason":"sample reason"},"requestedAt":"0001-01-01T00:00:00Z","requestedBy":{"email":"user1@gmail.com","firstName":"","id":"","lastName":""},"status":"ACTIVE"}`,
			createRes: r,
		},

		{
			name: "no preflight ID",
			give: `{
  "groupOptions": [
    {
      "id": "group1",
      "timing": {
        "durationSeconds": 3600
      }
    },
    {
      "id": "group2",
      "timing": {
        "durationSeconds": 1800
      }
    }
  ],
  "preflightId": "1234567890",
	"createTemplate": false,
  "reason": "Sample reason"
}`,
			wantCode:      http.StatusNotFound,
			mockCreateErr: accesssvc.ErrPreflightNotFound,
			wantBody:      `{"error":"preflight not found"}`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			mockAccess := mocks.NewMockAccessService(ctrl)

			mockAccess.EXPECT().CreateRequest(gomock.Any(), gomock.Any(), gomock.Any()).Return(&tc.createRes, tc.mockCreateErr).AnyTimes()
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

	for _, tc := range testcases {
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
		givenID           string
		mockGetRequestErr error
		// request body (Request type)
		mockGetRequest *access.RequestWithGroupsWithTargets

		// expected HTTP response code
		wantCode int
		// expected HTTP response body
		wantBody string
		// withRequestArgumentsResponse map[string]types.RequestArgument
	}

	testcases := []testcase{
		{
			name:           "requestor can see their own request",
			givenID:        `req_123`,
			wantCode:       http.StatusOK,
			mockGetRequest: &r,
			// canReview is false in the response
			wantBody: `{"accessGroups":[{"accessRule":{"timeConstraints":{"maxDurationSeconds":0}},"createdAt":"0001-01-01T00:00:00Z","id":"group1","requestId":"req_123","requestStatus":"ACTIVE","requestedBy":{"email":"abc@gmail.com","firstName":"abc","id":"123","lastName":"xyz"},"requestedTiming":{"durationSeconds":0},"status":"APPROVED","targets":[{"accessGroupId":"group1","fields":[{"fieldTitle":"","id":"123","value":"","valueLabel":""}],"id":"xyz","requestId":"req_123","requestedBy":{"email":"abc@gmail.com","firstName":"abc","id":"123","lastName":"xyz"},"status":"ACTIVE","targetGroupId":"aws","targetKind":{"icon":"","kind":"","name":"","publisher":""}}],"updatedAt":"0001-01-01T00:00:00Z"},{"accessRule":{"timeConstraints":{"maxDurationSeconds":0}},"createdAt":"0001-01-01T00:00:00Z","id":"group2","requestId":"123","requestStatus":"ACTIVE","requestedBy":{"email":"","firstName":"","id":"","lastName":""},"requestedTiming":{"durationSeconds":0},"status":"APPROVED","targets":[{"accessGroupId":"group2","fields":[{"fieldTitle":"","id":"123","value":"","valueLabel":""}],"id":"xyz","requestId":"req_123","requestedBy":{"email":"","firstName":"","id":"","lastName":""},"status":"ACTIVE","targetGroupId":"aws","targetKind":{"icon":"","kind":"","name":"","publisher":""}}],"updatedAt":"0001-01-01T00:00:00Z"}],"id":"req_123","purpose":{"reason":"sample reason"},"requestedAt":"0001-01-01T00:00:00Z","requestedBy":{"email":"user1@gmail.com","firstName":"","id":"","lastName":""},"status":"ACTIVE"}`,
		},
		{
			name:              "noRequestFound",
			givenID:           `wrongID`,
			wantCode:          http.StatusNotFound,
			mockGetRequestErr: ddb.ErrNoItems,
			wantBody:          `{"error":"item query returned no items"}`,
		},
		{
			name:              "not a requestor or reviewer",
			givenID:           `req_123`,
			mockGetRequest:    nil,
			wantCode:          http.StatusNotFound,
			mockGetRequestErr: ddb.ErrNoItems,
			wantBody:          `{"error":"item query returned no items"}`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {

			db := ddbmock.New(t)
			db.MockQueryWithErr(&storage.GetRequestWithGroupsWithTargetsForUserOrReviewer{Result: tc.mockGetRequest}, tc.mockGetRequestErr)

			a := API{DB: db}
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
		giveFilter     *string
		mockDBQuery    ddb.QueryBuilder
		mockDBQueryErr error
		// expected HTTP response code
		wantCode int
		// expected HTTP response body
		wantBody string
	}

	testcases := []testcase{

		{
			name:     "ok requestor",
			wantCode: http.StatusOK,
			mockDBQuery: &storage.ListRequestWithGroupsWithTargetsForUser{Result: []access.RequestWithGroupsWithTargets{{
				Request: r.Request,
				Groups:  r.Groups,
			}}},
			wantBody: `{"next":null,"requests":[{"accessGroups":[{"accessRule":{"timeConstraints":{"maxDurationSeconds":0}},"createdAt":"0001-01-01T00:00:00Z","id":"group1","requestId":"req_123","requestStatus":"ACTIVE","requestedBy":{"email":"abc@gmail.com","firstName":"abc","id":"123","lastName":"xyz"},"requestedTiming":{"durationSeconds":0},"status":"APPROVED","targets":[{"accessGroupId":"group1","fields":[{"fieldTitle":"","id":"123","value":"","valueLabel":""}],"id":"xyz","requestId":"req_123","requestedBy":{"email":"abc@gmail.com","firstName":"abc","id":"123","lastName":"xyz"},"status":"ACTIVE","targetGroupId":"aws","targetKind":{"icon":"","kind":"","name":"","publisher":""}}],"updatedAt":"0001-01-01T00:00:00Z"},{"accessRule":{"timeConstraints":{"maxDurationSeconds":0}},"createdAt":"0001-01-01T00:00:00Z","id":"group2","requestId":"123","requestStatus":"ACTIVE","requestedBy":{"email":"","firstName":"","id":"","lastName":""},"requestedTiming":{"durationSeconds":0},"status":"APPROVED","targets":[{"accessGroupId":"group2","fields":[{"fieldTitle":"","id":"123","value":"","valueLabel":""}],"id":"xyz","requestId":"req_123","requestedBy":{"email":"","firstName":"","id":"","lastName":""},"status":"ACTIVE","targetGroupId":"aws","targetKind":{"icon":"","kind":"","name":"","publisher":""}}],"updatedAt":"0001-01-01T00:00:00Z"}],"id":"req_123","purpose":{"reason":"sample reason"},"requestedAt":"0001-01-01T00:00:00Z","requestedBy":{"email":"user1@gmail.com","firstName":"","id":"","lastName":""},"status":"ACTIVE"}]}`,
		},
		{
			name:           "ok with no requests",
			wantCode:       http.StatusOK,
			mockDBQuery:    &storage.ListRequestWithGroupsWithTargetsForUser{Result: nil},
			mockDBQueryErr: nil,
			wantBody:       `{"next":null,"requests":[]}`,
		},
		{
			name:           "unhandled error",
			wantCode:       http.StatusInternalServerError,
			mockDBQuery:    &storage.ListRequestWithGroupsWithTargetsForUser{},
			mockDBQueryErr: errors.New("random error"),
			wantBody:       `{"error":"Internal Server Error"}`,
		},
		{
			name:       "with filter param PAST",
			giveFilter: aws.String("PAST"),
			wantCode:   http.StatusOK,
			mockDBQuery: &storage.ListRequestWithGroupsWithTargetsForUserAndPastUpcoming{
				Result: []access.RequestWithGroupsWithTargets{{
					Request: r.Request,
					Groups:  r.Groups,
				}}},
			wantBody: `{"next":null,"requests":[{"accessGroups":[{"accessRule":{"timeConstraints":{"maxDurationSeconds":0}},"createdAt":"0001-01-01T00:00:00Z","id":"group1","requestId":"req_123","requestStatus":"ACTIVE","requestedBy":{"email":"abc@gmail.com","firstName":"abc","id":"123","lastName":"xyz"},"requestedTiming":{"durationSeconds":0},"status":"APPROVED","targets":[{"accessGroupId":"group1","fields":[{"fieldTitle":"","id":"123","value":"","valueLabel":""}],"id":"xyz","requestId":"req_123","requestedBy":{"email":"abc@gmail.com","firstName":"abc","id":"123","lastName":"xyz"},"status":"ACTIVE","targetGroupId":"aws","targetKind":{"icon":"","kind":"","name":"","publisher":""}}],"updatedAt":"0001-01-01T00:00:00Z"},{"accessRule":{"timeConstraints":{"maxDurationSeconds":0}},"createdAt":"0001-01-01T00:00:00Z","id":"group2","requestId":"123","requestStatus":"ACTIVE","requestedBy":{"email":"","firstName":"","id":"","lastName":""},"requestedTiming":{"durationSeconds":0},"status":"APPROVED","targets":[{"accessGroupId":"group2","fields":[{"fieldTitle":"","id":"123","value":"","valueLabel":""}],"id":"xyz","requestId":"req_123","requestedBy":{"email":"","firstName":"","id":"","lastName":""},"status":"ACTIVE","targetGroupId":"aws","targetKind":{"icon":"","kind":"","name":"","publisher":""}}],"updatedAt":"0001-01-01T00:00:00Z"}],"id":"req_123","purpose":{"reason":"sample reason"},"requestedAt":"0001-01-01T00:00:00Z","requestedBy":{"email":"user1@gmail.com","firstName":"","id":"","lastName":""},"status":"ACTIVE"}]}`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			//will need to test both upcoming and current lists

			db := ddbmock.New(t)
			db.MockQueryWithErr(tc.mockDBQuery, tc.mockDBQueryErr)
			a := API{DB: db}
			handler := newTestServer(t, &a)
			var qp []string

			if tc.giveFilter != nil {
				qp = append(qp, "filter="+*tc.giveFilter)
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
			}

			if tc.wantBody != "" {
				assert.Equal(t, tc.wantBody, string(data))
			}
		})
	}

}

func TestRevokeRequest(t *testing.T) {
	type testcase struct {
		request              access.RequestWithGroupsWithTargets
		name                 string
		give                 string
		revokeSvcResp        *access.RequestWithGroupsWithTargets
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
			request:             r,
			wantCode:            http.StatusNotFound,
			withGetRequestError: ddb.ErrNoItems,
			wantBody:            `{"error":"request not found or you don't have access to it"}`,
		},
		{
			name:          "user can revoke their own grant",
			request:       r,
			withUID:       "user1",
			withUEmail:    "user1@gmail.com",
			wantCode:      http.StatusOK,
			revokeSvcResp: &access.RequestWithGroupsWithTargets{},
			withIsAdmin:   false,
		},
		{
			name:          "admin can revoke any request",
			revokeSvcResp: &access.RequestWithGroupsWithTargets{},
			withUID:       "admin",
			withUEmail:    "admin@mail.com",
			withIsAdmin:   true,
			wantCode:      http.StatusOK,
		},
		{
			name:                "user cant revoke other users request",
			withIsAdmin:         false,
			withUID:             "abcd",
			withUEmail:          "userinvalid@gmai.com",
			request:             r,
			wantCode:            http.StatusNotFound,
			withGetRequestError: ddb.ErrNoItems,
			wantBody:            `{"error":"request not found or you don't have access to it"}`,
		},
		{
			name:          "reviewer can revoke request",
			request:       r,
			revokeSvcResp: &access.RequestWithGroupsWithTargets{},
			withUID:       "user2",
			withUEmail:    "user2@mail.com",
			wantCode:      http.StatusOK,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {

			db := ddbmock.New(t)
			db.MockQueryWithErr(&storage.GetRequestWithGroupsWithTargets{Result: &tc.request}, tc.withGetRequestError)
			db.MockQueryWithErr(&storage.GetRequestReviewer{Result: &access.Reviewer{
				ReviewerID: tc.request.Request.RequestedBy.ID,
				RequestID:  tc.request.Request.ID,
			}}, tc.withGetReviewerError)
			ctrl := gomock.NewController(t)
			m := mocks.NewMockAccessService(ctrl)

			m.EXPECT().RevokeRequest(gomock.Any(), gomock.Any()).Return(tc.revokeSvcResp, tc.withRevokeGrantErr).AnyTimes()
			a := API{DB: db, Access: m}
			handler := newTestServer(t, &a, WithIsAdmin(tc.withIsAdmin), WithRequestUser(identity.User{ID: tc.withUID, Email: tc.withUEmail}))

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
		mockGetRequest            storage.GetRequestWithGroupsWithTargets
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
			mockGetRequest: storage.GetRequestWithGroupsWithTargets{
				ID:     "",
				Result: &r,
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
			name:     "empty event lists",
			wantCode: http.StatusOK,
			mockGetRequest: storage.GetRequestWithGroupsWithTargets{
				ID:     "",
				Result: &r,
			},
			mockGetRequestReviewer: storage.GetRequestReviewer{
				RequestID:  "1234",
				ReviewerID: "abcd",
				Result: &access.Reviewer{
					ReviewerID:    "abcd",
					RequestID:     "req_123",
					Notifications: access.Notifications{},
				},
			},
			mockListEvents: storage.ListRequestEvents{
				RequestID: "1234",
				Result:    []access.RequestEvent{},
			},
			apiUserID: "abcd",
			wantBody:  `{"events":[],"next":null}`,
		},
		{
			name:              "not found",
			wantCode:          http.StatusUnauthorized,
			mockGetRequestErr: ddb.ErrNoItems,

			wantBody: `{"error":"item query returned no items"}`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {

			db := ddbmock.New(t)
			db.MockQueryWithErr(&tc.mockGetRequest, tc.mockGetRequestErr)
			db.MockQueryWithErr(&tc.mockListEvents, tc.mockListEventsErr)
			db.MockQueryWithErr(&tc.mockGetRequestReviewer, tc.mockGetRequestReviewerErr)
			a := API{DB: db}
			handler := newTestServer(t, &a, WithRequestUser(identity.User{ID: tc.apiUserID}), WithIsAdmin(tc.apiUserIsAdmin))

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
			}

			if tc.wantBody != "" {
				assert.Equal(t, tc.wantBody, string(data))
			}
		})
	}

}
