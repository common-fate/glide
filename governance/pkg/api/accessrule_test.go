package api

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/common-fate/common-fate/governance/pkg/api/mocks"
	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/common-fate/common-fate/pkg/service/rulesvc"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/ddb"
	"github.com/common-fate/ddb/ddbmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGovListAccessRules(t *testing.T) {
	type testcase struct {
		name        string
		rules       []rule.AccessRule
		want        string
		mockListErr error
		wantCode    int
	}

	testcases := []testcase{
		{
			name:     "ok",
			wantCode: http.StatusOK,
			rules: []rule.AccessRule{
				{
					ID:          "rule1",
					Status:      rule.ACTIVE,
					Description: "string",
					Name:        "string",
					Groups:      []string{"string"},
					Target: rule.Target{
						ProviderID:          "string",
						BuiltInProviderType: "okta",
						With:                map[string]string{},
					},
					Approval: rule.Approval{
						Groups: []string{"a"},
						Users:  []string{"b"},
					},
				},
				{
					ID:          "rule2",
					Status:      rule.ACTIVE,
					Description: "string",
					Name:        "string",
					Groups:      []string{"string"},
					Target: rule.Target{
						ProviderID:          "string",
						BuiltInProviderType: "okta",
						With:                map[string]string{},
					},
				},
			},

			want: `{"accessRules":[{"approval":{"groups":["a"],"users":["b"]},"description":"string","groups":["string"],"id":"rule1","isCurrent":false,"metadata":{"createdAt":"0001-01-01T00:00:00Z","createdBy":"","updatedAt":"0001-01-01T00:00:00Z","updatedBy":""},"name":"string","status":"ACTIVE","target":{"provider":{"id":"string","type":"okta"},"with":{}},"timeConstraints":{"maxDurationSeconds":0},"version":""},{"approval":{"groups":[],"users":[]},"description":"string","groups":["string"],"id":"rule2","isCurrent":false,"metadata":{"createdAt":"0001-01-01T00:00:00Z","createdBy":"","updatedAt":"0001-01-01T00:00:00Z","updatedBy":""},"name":"string","status":"ACTIVE","target":{"provider":{"id":"string","type":"okta"},"with":{}},"timeConstraints":{"maxDurationSeconds":0},"version":""}],"next":null}`,
		},
		{
			name:        "no rules returns an empty list not an error",
			mockListErr: ddb.ErrNoItems,
			wantCode:    http.StatusOK,
			rules:       nil,

			want: `{"accessRules":[],"next":null}`,
		},
		{
			name:        "internal error",
			mockListErr: errors.New("internal error"),
			wantCode:    http.StatusInternalServerError,
			rules:       nil,

			want: `{"error":"Internal Server Error"}`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			db := ddbmock.New(t)
			db.MockQueryWithErr(&storage.ListCurrentAccessRules{Result: tc.rules}, tc.mockListErr)

			a := API{DB: db}
			handler := newTestServer(t, &a)

			req, err := http.NewRequest("GET", "/gov/v1/access-rules", nil)
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

			assert.Equal(t, tc.want, string(data))
		})
	}
}

func TestGovCreateAccessRule(t *testing.T) {
	type testcase struct {
		name          string
		give          string
		mockCreate    *rule.AccessRule
		mockCreateErr error
		wantCode      int
		wantBody      string
	}

	testcases := []testcase{
		{
			name: "ok",
			give: `{"target":{"providerId":"string","with":{}},"timeConstraints":{"maxDurationSeconds": 60},"groups":["string"],"name":"string","description":"string","approval":{"groups":[],"users":[]}}`,
			mockCreate: &rule.AccessRule{
				ID:          "rule1",
				Status:      rule.ACTIVE,
				Description: "string",
				Name:        "string",
				Groups:      []string{"string"},
				Target: rule.Target{
					ProviderID: "string",
					With:       map[string]string{},
				},
				Metadata: rule.AccessRuleMetadata{
					CreatedBy: "bot_governance_api",
				},
			},
			wantCode: http.StatusCreated,
			wantBody: `{"approval":{"groups":[],"users":[]},"description":"string","groups":["string"],"id":"rule1","isCurrent":false,"metadata":{"createdAt":"0001-01-01T00:00:00Z","createdBy":"bot_governance_api","updatedAt":"0001-01-01T00:00:00Z","updatedBy":""},"name":"string","status":"ACTIVE","target":{"provider":{"id":"string","type":""},"with":{}},"timeConstraints":{"maxDurationSeconds":0},"version":""}`,
		},
		{
			name:          "id already exists",
			give:          `{"target":{"providerId":"string","with":{}},"timeConstraints":{"maxDurationSeconds": 60},"groups":["string"],"name":"string","description":"string","approval":{"groups":[],"users":[]}}`,
			mockCreateErr: rulesvc.ErrRuleIdAlreadyExists,
			wantCode:      http.StatusBadRequest,
			wantBody:      `{"error":"access rule id already exists"}`,
		},
		{
			name:     "fail when rule doesn't meet maxduration req",
			give:     `{"target":{"providerId":"string","with":{}},"timeConstraints":{"maxDurationSeconds": 1},"groups":["string"],"name":"string","description":"string","approval":{"groups":[],"users":[]}}`,
			wantCode: http.StatusBadRequest,
			wantBody: `{"error":"request body has an error: doesn't match the schema: Error at \"/timeConstraints/maxDurationSeconds\": number must be at least 60"}`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mocks.NewMockAccessRuleService(ctrl)
			if (tc.mockCreate != nil) || (tc.mockCreateErr != nil) {
				m.EXPECT().CreateAccessRule(gomock.Any(), gomock.Any(), gomock.Any()).Return(tc.mockCreate, tc.mockCreateErr)
			}

			db := ddbmock.New(t)

			db.MockQuery(&storage.ListUsers{})
			db.MockQuery(&storage.ListGroups{})

			a := API{Rules: m, DB: db}
			handler := newTestServer(t, &a)

			req, err := http.NewRequest("POST", "/gov/v1/access-rules", strings.NewReader(tc.give))
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
