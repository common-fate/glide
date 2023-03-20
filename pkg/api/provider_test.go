package api

// import (
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/common-fate/provider-registry/pkg/types"
// 	"github.com/stretchr/testify/assert"

// 	"github.com/golang/mock/gomock"
// )

// func TestListProviders(t *testing.T) {

// 	type testcase struct {
// 		name          string
// 		mockCreate    *types.ListProvidersResponse
// 		mockCreateErr error
// 		wantCode      int
// 		wantBody      string
// 	}

// 	list := []types.Provider{
// 		{
// 			Id:   "cf-dev",
// 			Type: "aws-sso",
// 		},
// 	}

// 	emptyList := []types.Provider{}

// 	errorMsg := "internal server error"

// 	testcases := []testcase{
// 		{
// 			name:     "ok",
// 			wantCode: http.StatusOK,
// 			mockCreate: &types.ListProvidersResponse{
// 				JSON200:      &list,
// 				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
// 			},
// 			wantBody: `[{"id":"cf-dev","type":"aws-sso"}]`,
// 		},
// 		{
// 			name:     "empty list should return empty array []",
// 			wantCode: http.StatusOK,
// 			mockCreate: &types.ListProvidersResponse{
// 				JSON200:      &emptyList,
// 				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
// 			},
// 			wantBody: `[]`,
// 		},
// 		{
// 			name:     "JSON500 should return error message",
// 			wantCode: http.StatusInternalServerError,
// 			mockCreate: &types.ListProvidersResponse{
// 				HTTPResponse: &http.Response{StatusCode: http.StatusInternalServerError},
// 				JSON500: &struct {
// 					Error *string "json:\"error,omitempty\""
// 				}{
// 					Error: &errorMsg,
// 				},
// 			},
// 			wantBody: `{"error":"internal server error"}`,
// 		},
// 		{
// 			name:     "unhandled should return generic error message",
// 			wantCode: http.StatusInternalServerError,
// 			mockCreate: &types.ListProvidersResponse{
// 				HTTPResponse: &http.Response{StatusCode: http.StatusBadGateway},
// 			},
// 			wantBody: `{"error":"Internal Server Error"}`,
// 		},
// 	}

// 	for _, tc := range testcases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			t.Parallel()

// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			m := ahmocks.NewMockClientWithResponsesInterface(ctrl)
// 			m.EXPECT().ListProvidersWithResponse(gomock.Any(), gomock.Any()).Return(tc.mockCreate, tc.mockCreateErr)

// 			a := API{
// 				AccessHandlerClient: m,
// 			}

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
