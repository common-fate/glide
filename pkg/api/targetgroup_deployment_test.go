package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/common-fate/common-fate/pkg/api/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

/**
Help me stub out test cases for
TestListTargetGroupDeployments
TestCreateTargetGroupDeployment
TestGetTargetGroupDeployment
*/

func TestListTargetGroupDeployments(t *testing.T) {
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
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			mockAccess := mocks.NewMockAccessService(ctrl)
			mockAccess.EXPECT().CancelRequest(gomock.Any(), gomock.Any()).Return(tc.mockCancelErr).AnyTimes()
			a := API{Access: mockAccess}
			handler := newTestServer(t, &a)

			req, err := http.NewRequest("GET", "/api/v1/admin/access/cancel", strings.NewReader(tc.give))
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

func TestCreateTargetGroupDeployment(t *testing.T) {
	// @TODO:

}

func TestGetTargetGroupDeployment(t *testing.T) {
	// @TODO:

}
