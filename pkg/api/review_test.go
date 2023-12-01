package api

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/api/mocks"
	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/common-fate/common-fate/pkg/service/accesssvc"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/ddb/ddbmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestReviewRequest(t *testing.T) {
	type testcase struct {
		name              string
		give              string
		requestErr        error // error from fetching the request
		reviewersErr      error // error from fetching reviewers
		addReviewResult   *accesssvc.AddReviewResult
		addReviewErr      error
		withTestUser      *identity.User
		wantAddReviewOpts accesssvc.AddReviewOpts
		wantCode          int
		wantBody          string
	}
	overrideTime, err := time.Parse("2006-01-02T15:04:05.999Z", "2020-01-01T16:20:10Z")
	if err != nil {
		t.Fatal(err)
	}
	testcases := []testcase{
		{
			name: "ok",
			give: `{"decision": "APPROVED"}`,
			addReviewResult: &accesssvc.AddReviewResult{
				// fill the struct a little bit to verify it is included in the HTTP response
				Request: access.Request{
					ID: "test",
				},
			},
			wantAddReviewOpts: accesssvc.AddReviewOpts{Decision: access.DecisionApproved},
			wantCode:          http.StatusCreated,
			wantBody:          `{"request":{"accessRuleId":"","accessRuleVersion":"","id":"test","requestedAt":"0001-01-01T00:00:00Z","requestor":"","status":"","timing":{"durationSeconds":0},"updatedAt":"0001-01-01T00:00:00Z"}}`,
		},
		{
			name: "ok with override timing",
			give: `{"decision": "APPROVED","overrideTiming":{"durationSeconds":3600,"startTime":"2020-01-01T16:20:10Z"}}`,
			addReviewResult: &accesssvc.AddReviewResult{
				// fill the struct a little bit to verify it is included in the HTTP responses
				Request: access.Request{
					ID: "test",
					OverrideTiming: &access.Timing{
						Duration:  time.Second * 3600,
						StartTime: &overrideTime,
					},
				},
			},
			wantAddReviewOpts: accesssvc.AddReviewOpts{Decision: access.DecisionApproved, OverrideTiming: &access.Timing{
				Duration:  time.Second * 3600,
				StartTime: &overrideTime,
			}},
			wantCode: http.StatusCreated,
			wantBody: `{"request":{"accessRuleId":"","accessRuleVersion":"","id":"test","requestedAt":"0001-01-01T00:00:00Z","requestor":"","status":"","timing":{"durationSeconds":3600,"startTime":"2020-01-01T16:20:10Z"},"updatedAt":"0001-01-01T00:00:00Z"}}`,
		},
		{
			name:              "not authorized",
			give:              `{"decision": "APPROVED"}`,
			addReviewErr:      accesssvc.ErrUserNotAuthorized,
			wantCode:          http.StatusUnauthorized,
			wantAddReviewOpts: accesssvc.AddReviewOpts{Decision: access.DecisionApproved},
			wantBody:          `{"error":"you are not a reviewer of this request"}`,
		},
		{
			name:              "review fetching error",
			give:              `{"decision": "APPROVED"}`,
			reviewersErr:      errors.New("error"),
			wantCode:          http.StatusInternalServerError,
			wantAddReviewOpts: accesssvc.AddReviewOpts{Decision: access.DecisionApproved},
			wantBody:          `{"error":"Internal Server Error"}`,
		},
		{
			name:              "admin can approve",
			wantAddReviewOpts: accesssvc.AddReviewOpts{Decision: access.DecisionApproved, ReviewerIsAdmin: true},
			withTestUser:      &identity.User{Groups: []string{"testAdmin"}},
			give:              `{"decision": "APPROVED"}`,
			addReviewErr:      accesssvc.ErrUserNotAuthorized,
			wantCode:          http.StatusUnauthorized,
			wantBody:          `{"error":"you are not a reviewer of this request"}`,
		},
		{
			name:              "multiple decision race condition",
			wantAddReviewOpts: accesssvc.AddReviewOpts{Decision: access.DecisionApproved, ReviewerIsAdmin: true},
			withTestUser:      &identity.User{Groups: []string{"testAdmin"}},
			give:              `{"decision": "APPROVED"}`,
			addReviewErr:      accesssvc.ErrUserNotAuthorized,
			wantCode:          http.StatusUnauthorized,
			wantBody:          `{"error":"you are not a reviewer of this request"}`,
		},
	}

	for i := range testcases {
		tc := testcases[i]
		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			mockAccess := mocks.NewMockAccessService(ctrl)

			db := ddbmock.New(t)

			mockAccess.EXPECT().AddReviewAndGrantAccess(gomock.Any(), tc.wantAddReviewOpts).Return(tc.addReviewResult, tc.addReviewErr).AnyTimes()
			db.MockQueryWithErr(&storage.ListRequestReviewers{}, tc.reviewersErr)
			db.MockQueryWithErr(&storage.GetRequest{Result: &access.Request{}}, tc.requestErr)
			db.MockQueryWithErr(&storage.GetAccessRuleCurrent{Result: &rule.AccessRule{}}, nil)

			a := API{Access: mockAccess, DB: db, AdminGroup: "testAdmin"}
			opts := []func(*testOptions){}
			if tc.withTestUser != nil {
				opts = append(opts, withRequestUser(*tc.withTestUser))
			}
			handler := newTestServer(t, &a, opts...)

			req, err := http.NewRequest("POST", "/api/v1/requests/abcd/review", strings.NewReader(tc.give))
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
