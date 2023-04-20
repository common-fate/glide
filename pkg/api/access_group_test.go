package api

import (
	"net/http"
	"testing"

	"github.com/common-fate/common-fate/pkg/access"
)

func TestListAccessGroups(t *testing.T) {
	type testcase struct {
		name         string
		accessGroups []access.Group
		wantCode     int
		wantBody     string
	}

	testcases := []testcase{
		{
			name:     "ok",
			wantCode: http.StatusOK,
			// accessGroups: []access.Group{
			// 	{
			// 		ID: "123",
			// 		AccessRule: rule.AccessRule{
			// 			ID: "abc",
			// 		},

			// 		Request: "test",
			// 		Status:  access.APPROVED,
			// 	},
			// 	{
			// 		ID: "456",
			// 		AccessRule: rule.AccessRule{
			// 			ID: "abc",
			// 		},

			// 		Request: "test",
			// 		Status:  access.APPROVED,
			// 	},
			// },
			wantBody: `{"groups":[{"grants":[],"id":"123","overrideTiming":{"durationSeconds":0},"request":"test","status":"APPROVED","time":{"maxDurationSeconds":0},"with":[{}]},{"grants":[],"id":"456","overrideTiming":{"durationSeconds":0},"request":"test","status":"APPROVED","time":{"maxDurationSeconds":0},"with":[{}]}]}`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// t.Parallel()
			// db := ddbmock.New(t)
			// db.MockQuery(&storage.ListAccessGroups{Result: tc.accessGroups})

			// a := API{DB: db}
			// handler := newTestServer(t, &a)

			// req, err := http.NewRequest("GET", "/api/v1/requests/test/groups", nil)
			// if err != nil {
			// 	t.Fatal(err)
			// }
			// req.Header.Add("Content-Type", "application/json")

			// rr := httptest.NewRecorder()

			// handler.ServeHTTP(rr, req)

			// assert.Equal(t, tc.wantCode, rr.Code)

			// data, err := io.ReadAll(rr.Body)
			// if err != nil {
			// 	t.Fatal(err)
			// }

			// assert.Equal(t, tc.wantBody, string(data))
		})
	}
}

func TestGetAccessGroup(t *testing.T) {
	type testcase struct {
		name        string
		idpErr      error
		accessGroup *access.Group
		wantCode    int
		wantBody    string
	}

	testcases := []testcase{
		// {
		// 	name:     "ok",
		// 	wantCode: http.StatusOK,
		// 	accessGroup: &access.Group{

		// 		ID: "123",
		// 		AccessRule: rule.AccessRule{
		// 			ID: "abc",
		// 		},

		// 		Request: "test",
		// 		Status:  access.APPROVED,
		// 	},
		// 	wantBody: `{"grants":[],"id":"123","overrideTiming":{"durationSeconds":0},"request":"test","status":"APPROVED","time":{"maxDurationSeconds":0},"with":[{"foo":"bar"}]}`,
		// },
		// {
		// 	name:     "group not found",
		// 	wantCode: http.StatusNotFound,
		// 	idpErr:   ddb.ErrNoItems,

		// 	wantBody: `{"error":"item query returned no items"}`,
		// },
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// t.Parallel()
			// db := ddbmock.New(t)
			// db.MockQueryWithErr(&storage.GetAccessGroups{Result: tc.accessGroup}, tc.idpErr)

			// a := API{DB: db}
			// handler := newTestServer(t, &a)

			// req, err := http.NewRequest("GET", "/api/v1/requests/test/groups/123", nil)
			// if err != nil {
			// 	t.Fatal(err)
			// }
			// req.Header.Add("Content-Type", "application/json")

			// rr := httptest.NewRecorder()

			// handler.ServeHTTP(rr, req)

			// assert.Equal(t, tc.wantCode, rr.Code)

			// data, err := io.ReadAll(rr.Body)
			// if err != nil {
			// 	t.Fatal(err)
			// }

			// assert.Equal(t, tc.wantBody, string(data))
		})
	}
}
