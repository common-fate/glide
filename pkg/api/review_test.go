package api

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/api/mocks"
	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/common-fate/common-fate/pkg/service/accesssvc"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb/ddbmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestReviewRequest(t *testing.T) {
	type testcase struct {
		name         string
		withTestUser *identity.User
		// give contains post request body
		give         string
		addReviewErr error
		wantCode     int
		wantBody     string
	}

	testcases := []testcase{
		{
			name:         "ok",
			give:         `{"decision": "APPROVED"}`,
			wantCode:     http.StatusNoContent,
			withTestUser: &identity.User{Groups: []string{"testAdmin"}},
			wantBody:     "",
		},
		{
			name:         "ok with override timing",
			give:         `{"decision": "APPROVED","overrideTiming":{"durationSeconds":3600,"startTime":"2020-01-01T16:20:10Z"}}`,
			wantCode:     http.StatusBadRequest,
			withTestUser: &identity.User{Groups: []string{"testAdmin"}},
			addReviewErr: accesssvc.ErrGroupCannotBeApprovedBecauseItWillOverlapExistingGrants,
			wantBody:     `{"error":"this group has grants which overlap with existing grants"}`,
		},
		{
			name:         "not authorized",
			give:         `{"decision": "APPROVED"}`,
			withTestUser: &identity.User{Groups: []string{"invalid-group"}},
			addReviewErr: accesssvc.ErrAccesGroupNotFoundOrNoAccessToReview,
			wantCode:     http.StatusUnauthorized,
			wantBody:     `{"error":"this access group doesn't exist or you don't have access to review it"}`,
		},
		{
			name:         "unhandled error",
			give:         `{"decision": "APPROVED"}`,
			wantCode:     http.StatusInternalServerError,
			withTestUser: &identity.User{Groups: []string{"testAdmin"}},
			addReviewErr: errors.New("Internal Server Error"),
			wantBody:     `{"error":"Internal Server Error"}`,
		},
		{
			name:         "invalid request body",
			give:         `{"invalid": "APPROVED"}`,
			wantCode:     http.StatusBadRequest,
			withTestUser: &identity.User{Groups: []string{"testAdmin"}},
			addReviewErr: errors.New("Internal Server Error"),
			wantBody:     `{"error":"request body has an error: doesn't match the schema: Error at \"/decision\": property \"decision\" is missing"}`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			mockAccess := mocks.NewMockAccessService(ctrl)

			db := ddbmock.New(t)

			db.MockQueryWithErr(&storage.GetRequestGroupWithTargetsForReviewer{
				RequestID:  "abcd",
				GroupID:    "abcdef",
				ReviewerID: "asdf",
				Result: &access.GroupWithTargets{
					Group: access.Group{
						RequestID:        "abcd",
						RequestStatus:    types.PENDING,
						RequestReviewers: nil,
						GroupReviewers:   tc.withTestUser.Groups,
					},
					Targets: []access.GroupTarget{{
						RequestReviewers: tc.withTestUser.Groups,
					}},
				},
			}, nil)

			mockAccess.EXPECT().Review(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(tc.addReviewErr).AnyTimes()

			a := API{Access: mockAccess, DB: db, AdminGroup: "testAdmin"}
			opts := []func(*testOptions){}
			if tc.withTestUser != nil {
				opts = append(opts, WithRequestUser(*tc.withTestUser))
			}

			handler := newTestServer(t, &a, opts...)

			req, err := http.NewRequest("POST", "/api/v1/requests/abcd/review/abcdef", strings.NewReader(tc.give))
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
