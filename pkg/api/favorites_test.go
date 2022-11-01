package api

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/common-fate/ddb/ddbmock"
	"github.com/common-fate/granted-approvals/pkg/access"
	"github.com/common-fate/granted-approvals/pkg/api/mocks"
	"github.com/common-fate/granted-approvals/pkg/storage"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUserCreateFavorite(t *testing.T) {
	type testcase struct {
		name          string
		give          string
		mockCreate    *access.Favorite
		mockCreateErr error
		wantCode      int
		wantBody      string
	}

	testcases := []testcase{
		{
			name: "ok",
			give: `{"name":"test name","timing":{"durationSeconds": 10}, "accessRuleId": "rul_123"}`,
			mockCreate: &access.Favorite{
				ID:     "rqf_123",
				UserID: "usr_123",
				Name:   "test name",
				RequestedTiming: access.Timing{
					Duration: time.Second * 10,
				},
				Rule: "rul_123",
			},
			wantCode: http.StatusCreated,
			wantBody: `{"id":"rqf_123","name":"test name","timing":{"durationSeconds":10},"with":null}`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockAccess := mocks.NewMockAccessService(ctrl)
			mockAccess.EXPECT().CreateFavorite(gomock.Any(), gomock.Any(), gomock.Any()).Return(tc.mockCreate, tc.mockCreateErr).AnyTimes()
			a := API{Access: mockAccess}
			handler := newTestServer(t, &a)

			req, err := http.NewRequest("POST", "/api/v1/favorites", strings.NewReader(tc.give))
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

func TestUserGetFavorite(t *testing.T) {

	type testcase struct {
		name                string
		givenID             string
		mockGetFavoritetErr error
		mockGetFavorite     *access.Favorite
		// expected HTTP response code
		wantCode int
		// expected HTTP response body
		wantBody string
	}

	testcases := []testcase{
		{
			name:     "ok",
			givenID:  `rqf_123`,
			wantCode: http.StatusOK,
			mockGetFavorite: &access.Favorite{
				ID:     "rqf_123",
				UserID: "usr_123",
				Name:   "test name",
				RequestedTiming: access.Timing{
					Duration: time.Second * 10,
				},
				Rule: "rul_123",
				Data: access.RequestData{
					Reason: aws.String("a reason"),
				},
				With: []map[string][]string{
					{
						"accountId": {
							"abcd", "efgh",
						},
					},
				},
			},

			wantBody: `{"id":"rqf_123","name":"test name","reason":"a reason","timing":{"durationSeconds":10},"with":[{"accountId":["abcd","efgh"]}]}`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			db := ddbmock.New(t)
			db.MockQueryWithErr(&storage.GetFavoriteForUser{Result: tc.mockGetFavorite}, tc.mockGetFavoritetErr)
			a := API{DB: db}
			handler := newTestServer(t, &a)

			req, err := http.NewRequest("GET", "/api/v1/favorites/"+tc.givenID, strings.NewReader(""))
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

func TestUserListFavorites(t *testing.T) {

	type testcase struct {
		name           string
		mockFavorites  []access.Favorite
		mockDBQueryErr error
		// expected HTTP response code
		wantCode int
		// expected HTTP response body
		wantBody string
	}

	testcases := []testcase{
		{
			name:     "ok",
			wantCode: http.StatusOK,
			mockFavorites: []access.Favorite{
				{
					ID:     "rqf_123",
					UserID: "usr_123",
					Name:   "test name",
					RequestedTiming: access.Timing{
						Duration: time.Second * 10,
					},
					Rule: "rul_123",
					Data: access.RequestData{
						Reason: aws.String("a reason"),
					},
					With: []map[string][]string{
						{
							"accountId": {
								"abcd", "efgh",
							},
						},
					},
				},
			},

			wantBody: `{"favorites":[{"id":"rqf_123","name":"test name","ruleId":"rul_123"}],"next":null}`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			db := ddbmock.New(t)
			db.MockQueryWithErr(&storage.ListFavoritesForUser{Result: tc.mockFavorites}, tc.mockDBQueryErr)
			a := API{DB: db}
			handler := newTestServer(t, &a)
			req, err := http.NewRequest("GET", "/api/v1/favorites", strings.NewReader(""))
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
