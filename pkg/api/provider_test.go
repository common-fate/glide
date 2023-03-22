package api

// import (
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/common-fate/common-fate/pkg/storage"
// 	"github.com/common-fate/common-fate/pkg/target"
// 	"github.com/common-fate/ddb/ddbmock"
// 	"github.com/stretchr/testify/assert"

// 	"github.com/golang/mock/gomock"
// )

// func TestListProviders(t *testing.T) {

// 	type testcase struct {
// 		name     string
// 		listErr  error
// 		wantList []target.Group
// 		wantCode int
// 		wantBody string
// 	}

// 	testcases := []testcase{
// 		{
// 			name:     "ok",
// 			wantCode: http.StatusOK,
// 			wantList: []target.Group{{ID: "123"}},

// 			wantBody: `[{"id":"cf-dev","type":"aws-sso"}]`,
// 		},
// 		{
// 			name:     "empty list should return empty array []",
// 			wantCode: http.StatusOK,
// 			wantList: []target.Group{},
// 			wantBody: `[]`,
// 		},
// 		{
// 			name:     "JSON500 should return error message",
// 			wantCode: http.StatusInternalServerError,

// 			wantBody: `{"error":"internal server error"}`,
// 		},
// 		{
// 			name:     "unhandled should return generic error message",
// 			wantCode: http.StatusInternalServerError,

// 			wantBody: `{"error":"Internal Server Error"}`,
// 		},
// 	}

// 	for _, tc := range testcases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			t.Parallel()

// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()
// 			db := ddbmock.New(t)
// 			db.MockQueryWithErr(&storage.ListTargetGroups{Result: tc.wantList}, tc.listErr)

// 			a := API{DB: db}

// 			handler := newTestServer(t, &a)

// 			req, err := http.NewRequest("GET", "/api/v1/admin/providers", nil)
// 			if err != nil {
// 				t.Fatal(err)
// 			}

// 			req.Header.Add("Content-Type", "application/json")

// 			rr := httptest.NewRecorder()

// 			handler.ServeHTTP(rr, req)

// 			assert.Equal(t, tc.wantCode, rr.Code)

// 			assert.Equal(t, tc.wantBody, rr.Body.String())

// 		})
// 	}

// }
