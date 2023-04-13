package api

// import (
// 	"errors"
// 	"io"
// 	"net/http"
// 	"net/http/httptest"
// 	"strings"
// 	"testing"

// 	"github.com/common-fate/common-fate/pkg/api/mocks"
// 	"github.com/common-fate/common-fate/pkg/cache"
// 	"github.com/common-fate/common-fate/pkg/identity"
// 	"github.com/common-fate/common-fate/pkg/rule"
// 	"github.com/common-fate/common-fate/pkg/service/rulesvc"
// 	"github.com/common-fate/common-fate/pkg/storage"
// 	"github.com/common-fate/common-fate/pkg/types"
// 	"github.com/common-fate/ddb"
// 	"github.com/common-fate/ddb/ddbmock"
// 	"github.com/golang/mock/gomock"
// 	"github.com/stretchr/testify/assert"
// )

// func TestAdminCreateAccessRule(t *testing.T) {
// 	type testcase struct {
// 		name          string
// 		give          string
// 		mockCreate    *rule.AccessRule
// 		mockCreateErr error

// 		//idpUser  *types.User
// 		wantCode int
// 		wantBody string
// 	}

// 	testcases := []testcase{
// 		{
// 			name: "ok",
// 			give: `{"target":{"providerId":"string","with":{}},"timeConstraints":{"maxDurationSeconds": 60},"groups":["string"],"name":"string","description":"string","approval":{"groups":[],"users":[]}}`,
// 			mockCreate: &rule.AccessRule{
// 				ID:     "rule1",
// 				Status: rule.ACTIVE,

// 				Description: "string",
// 				Name:        "string",
// 				Groups:      []string{"string"},
// 				Target: rule.Target{
// 					TargetGroupID: "string",
// 					With:          map[string]string{},
// 				},
// 			},
// 			wantCode: http.StatusCreated,
// 			// idpUser: &types.User{
// 			// 	Id:    "userid",
// 			// 	Email: "test@test.com",
// 			// },

// 			wantBody: `{"approval":{"groups":[],"users":[]},"description":"string","groups":["string"],"id":"rule1","isCurrent":false,"metadata":{"createdAt":"0001-01-01T00:00:00Z","createdBy":"","updatedAt":"0001-01-01T00:00:00Z","updatedBy":""},"name":"string","status":"ACTIVE","target":{"provider":{"id":"string","type":""},"with":{}},"timeConstraints":{"maxDurationSeconds":0},"version":""}`,
// 		},
// 		{
// 			name:          "id already exists",
// 			give:          `{"target":{"providerId":"string","with":{}},"timeConstraints":{"maxDurationSeconds": 60},"groups":["string"],"name":"string","description":"string","approval":{"groups":[],"users":[]}}`,
// 			mockCreateErr: rulesvc.ErrRuleIdAlreadyExists,
// 			wantCode:      http.StatusBadRequest,
// 			wantBody:      `{"error":"access rule id already exists"}`,
// 		},
// 		{
// 			name:     "fail when rule doesn't meet maxduration req",
// 			give:     `{"target":{"providerId":"string","with":{}},"timeConstraints":{"maxDurationSeconds": 1},"groups":["string"],"name":"string","description":"string","approval":{"groups":[],"users":[]}}`,
// 			wantCode: http.StatusBadRequest,
// 			wantBody: `{"error":"request body has an error: doesn't match the schema: Error at \"/timeConstraints/maxDurationSeconds\": number must be at least 60"}`,
// 		},
// 	}

// 	for _, tc := range testcases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			t.Parallel()
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			m := mocks.NewMockAccessRuleService(ctrl)
// 			if (tc.mockCreate != nil) || (tc.mockCreateErr != nil) {
// 				m.EXPECT().CreateAccessRule(gomock.Any(), gomock.Any(), gomock.Any()).Return(tc.mockCreate, tc.mockCreateErr)
// 			}

// 			a := API{Rules: m}
// 			handler := newTestServer(t, &a)

// 			req, err := http.NewRequest("POST", "/api/v1/admin/access-rules", strings.NewReader(tc.give))
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			req.Header.Add("Content-Type", "application/json")
// 			rr := httptest.NewRecorder()

// 			handler.ServeHTTP(rr, req)

// 			assert.Equal(t, tc.wantCode, rr.Code)

// 			data, err := io.ReadAll(rr.Body)
// 			if err != nil {
// 				t.Fatal(err)
// 			}

// 			assert.Equal(t, tc.wantBody, string(data))
// 		})
// 	}

// }

// func TestAdminUpdateAccessRule(t *testing.T) {
// 	type testcase struct {
// 		name          string
// 		give          string
// 		mockCreate    *rule.AccessRule
// 		mockCreateErr error

// 		//idpUser  *types.User
// 		wantCode int
// 		wantBody string
// 		wantErr  string
// 	}

// 	testcases := []testcase{
// 		{
// 			name: "ok",
// 			give: `{"target":{"providerId":"string","with":{}},"approval":{"users":["a6936de0-633e-400b-8d36-5d3f47e1356e","629d4ea4-686c-4581-b778-ec083375523b"],"groups":[]},"name":"Productions","description":"Production access ","timeConstraints":{"maxDurationSeconds":3600},"groups":["common_fate_administrators"]}`,
// 			mockCreate: &rule.AccessRule{
// 				ID:          "rule1",
// 				Status:      rule.ACTIVE,
// 				Description: "Production access ",
// 				Name:        "Productions",
// 				Groups:      []string{"common_fate_administrators"},

// 				//target is not updated by this operation
// 				Target: rule.Target{
// 					TargetGroupID: "string",
// 					With:          map[string]string{},
// 				},
// 				Approval: rule.Approval{
// 					Groups: []string{},
// 					Users:  []string{"a6936de0-633e-400b-8d36-5d3f47e1356e", "629d4ea4-686c-4581-b778-ec083375523b"},
// 				},
// 				TimeConstraints: types.TimeConstraints{
// 					MaxDurationSeconds: 3600,
// 				},
// 			},
// 			wantCode: http.StatusAccepted,
// 			wantBody: `{"approval":{"groups":[],"users":["a6936de0-633e-400b-8d36-5d3f47e1356e","629d4ea4-686c-4581-b778-ec083375523b"]},"description":"Production access ","groups":["common_fate_administrators"],"id":"rule1","isCurrent":false,"metadata":{"createdAt":"0001-01-01T00:00:00Z","createdBy":"","updatedAt":"0001-01-01T00:00:00Z","updatedBy":""},"name":"Productions","status":"ACTIVE","target":{"provider":{"id":"string","type":""},"with":{}},"timeConstraints":{"maxDurationSeconds":3600},"version":"abcd"}`,
// 		},
// 		{
// 			name:     "malformed",
// 			give:     `malformed json input`,
// 			wantCode: http.StatusBadRequest,
// 			wantErr:  "{\"error\":\"request body has an error: failed to decode request body: invalid character 'm' looking for beginning of value\"}",
// 		},
// 	}

// 	for _, tc := range testcases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			t.Parallel()
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			m := mocks.NewMockAccessRuleService(ctrl)
// 			if tc.mockCreate != nil {
// 				m.EXPECT().UpdateRule(gomock.Any(), gomock.Any()).Return(tc.mockCreate, tc.mockCreateErr)
// 			}
// 			db := ddbmock.New(t)
// 			db.MockQuery(&storage.GetAccessRuleCurrent{Result: tc.mockCreate})
// 			a := API{Rules: m, DB: db}
// 			handler := newTestServer(t, &a)

// 			req, err := http.NewRequest("PUT", "/api/v1/admin/access-rules/"+"rule1", strings.NewReader(tc.give))
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			req.Header.Add("Content-Type", "application/json")
// 			rr := httptest.NewRecorder()

// 			handler.ServeHTTP(rr, req)

// 			assert.Equal(t, tc.wantCode, rr.Code)

// 			data, err := io.ReadAll(rr.Body)

// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			if tc.wantErr != "" {
// 				assert.Equal(t, tc.wantErr, string(data))
// 				return
// 			}
// 			if tc.wantBody != "" {
// 				assert.Equal(t, tc.wantBody, string(data))
// 			}
// 		})
// 	}
// }

// func TestAdminListAccessRules(t *testing.T) {
// 	type testcase struct {
// 		name string

// 		rules       []rule.AccessRule
// 		want        string
// 		mockListErr error
// 		wantCode    int
// 	}

// 	testcases := []testcase{
// 		{
// 			name: "ok",
// 			//mockListErr: types.ErrNoGroupsFound,
// 			wantCode: http.StatusOK,
// 			rules: []rule.AccessRule{
// 				{
// 					ID:          "rule1",
// 					Status:      rule.ACTIVE,
// 					Description: "string",
// 					Name:        "string",
// 					Groups:      []string{"string"},
// 					Target: rule.Target{
// 						TargetGroupID: "string",
// 						With:          map[string]string{},
// 					},
// 					// This should not be included in the response for users
// 					Approval: rule.Approval{
// 						Groups: []string{"a"},
// 						Users:  []string{"b"},
// 					},
// 				},
// 				{
// 					ID:          "rule2",
// 					Status:      rule.ACTIVE,
// 					Description: "string",
// 					Name:        "string",
// 					Groups:      []string{"string"},
// 					Target: rule.Target{
// 						TargetGroupID: "string",
// 						With:          map[string]string{},
// 					},
// 				},
// 			},

// 			want: `{"accessRules":[{"approval":{"groups":["a"],"users":["b"]},"description":"string","groups":["string"],"id":"rule1","isCurrent":false,"metadata":{"createdAt":"0001-01-01T00:00:00Z","createdBy":"","updatedAt":"0001-01-01T00:00:00Z","updatedBy":""},"name":"string","status":"ACTIVE","target":{"provider":{"id":"string","type":"okta"},"with":{}},"timeConstraints":{"maxDurationSeconds":0},"version":""},{"approval":{"groups":[],"users":[]},"description":"string","groups":["string"],"id":"rule2","isCurrent":false,"metadata":{"createdAt":"0001-01-01T00:00:00Z","createdBy":"","updatedAt":"0001-01-01T00:00:00Z","updatedBy":""},"name":"string","status":"ACTIVE","target":{"provider":{"id":"string","type":"okta"},"with":{}},"timeConstraints":{"maxDurationSeconds":0},"version":""}],"next":null}`,
// 		},
// 		{
// 			name:        "no rules returns an empty list not an error",
// 			mockListErr: ddb.ErrNoItems,
// 			wantCode:    http.StatusOK,
// 			rules:       nil,

// 			want: `{"accessRules":[],"next":null}`,
// 		},
// 		{
// 			name:        "internal error",
// 			mockListErr: errors.New("internal error"),
// 			wantCode:    http.StatusInternalServerError,
// 			rules:       nil,

// 			want: `{"error":"Internal Server Error"}`,
// 		},
// 	}

// 	for _, tc := range testcases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			t.Parallel()
// 			db := ddbmock.New(t)
// 			db.MockQueryWithErr(&storage.ListCurrentAccessRules{Result: tc.rules}, tc.mockListErr)

// 			a := API{DB: db}
// 			handler := newTestServer(t, &a)

// 			req, err := http.NewRequest("GET", "/api/v1/admin/access-rules", nil)
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			req.Header.Add("Content-Type", "application/json")
// 			rr := httptest.NewRecorder()

// 			handler.ServeHTTP(rr, req)

// 			assert.Equal(t, tc.wantCode, rr.Code)

// 			data, err := io.ReadAll(rr.Body)
// 			if err != nil {
// 				t.Fatal(err)
// 			}

// 			assert.Equal(t, tc.want, string(data))
// 		})
// 	}
// }

// func TestUserListAccessRules(t *testing.T) {
// 	type testcase struct {
// 		name string

// 		rules        []rule.AccessRule
// 		mockRulesErr error
// 		want         string
// 		wantCode     int
// 	}

// 	testcases := []testcase{
// 		{
// 			name:     "ok",
// 			wantCode: http.StatusOK,
// 			rules: []rule.AccessRule{
// 				{
// 					ID:          "rule1",
// 					Status:      rule.ACTIVE,
// 					Description: "string",
// 					Name:        "string",
// 					Groups:      []string{"string"},
// 					Target: rule.Target{
// 						TargetGroupID: "string",
// 						With:          map[string]string{},
// 					},
// 					// This should not be included in the response for users
// 					Approval: rule.Approval{
// 						Groups: []string{"a"},
// 						Users:  []string{"b"},
// 					},
// 				},
// 				{
// 					ID:          "rule2",
// 					Status:      rule.ACTIVE,
// 					Description: "string",
// 					Name:        "string",
// 					Groups:      []string{"string"},
// 					Target: rule.Target{
// 						TargetGroupID: "string",
// 						With:          map[string]string{},
// 					},
// 				},
// 			},

// 			want: `{"accessRules":[{"description":"string","id":"rule1","isCurrent":false,"name":"string","target":{"provider":{"id":"string","type":"okta"}},"timeConstraints":{"maxDurationSeconds":0},"version":""},{"description":"string","id":"rule2","isCurrent":false,"name":"string","target":{"provider":{"id":"string","type":"okta"}},"timeConstraints":{"maxDurationSeconds":0},"version":""}],"next":null}`,
// 		},
// 		{
// 			name:         "no rules found",
// 			mockRulesErr: ddb.ErrNoItems,
// 			wantCode:     http.StatusOK,
// 			want:         `{"accessRules":[],"next":null}`,
// 			rules:        []rule.AccessRule{},
// 		},
// 		{
// 			name:         "error fetching rules",
// 			mockRulesErr: errors.New("some error"),
// 			wantCode:     http.StatusInternalServerError,
// 			want:         `{"error":"Internal Server Error"}`,
// 			rules:        nil,
// 		},
// 	}

// 	for _, tc := range testcases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			t.Parallel()

// 			db := ddbmock.New(t)
// 			db.MockQueryWithErr(&storage.ListAccessRulesForStatus{Result: tc.rules}, tc.mockRulesErr)
// 			a := API{DB: db}
// 			handler := newTestServer(t, &a)

// 			req, err := http.NewRequest("GET", "/api/v1/access-rules", nil)
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			req.Header.Add("Content-Type", "application/json")
// 			rr := httptest.NewRecorder()

// 			handler.ServeHTTP(rr, req)

// 			assert.Equal(t, tc.wantCode, rr.Code)

// 			data, err := io.ReadAll(rr.Body)
// 			if err != nil {
// 				t.Fatal(err)
// 			}

// 			assert.Equal(t, tc.want, string(data))
// 		})
// 	}
// }
// func TestUserGetAccessRuleApprovals(t *testing.T) {
// 	type testcase struct {
// 		name                         string
// 		giveRuleID                   string
// 		mockGetRuleResponse          *rule.GetAccessRuleResponse
// 		mockGetRuleErr               error
// 		mockGetAccessRuleVersion     *rule.AccessRule
// 		withRequestArgumentsResponse map[string]types.RequestArgument
// 		want                         string
// 		wantCode                     int
// 	}

// 	testcases := []testcase{
// 		{
// 			name:       "ok",
// 			giveRuleID: "abcd",
// 			mockGetRuleResponse: &rule.GetAccessRuleResponse{
// 				Rule: &rule.AccessRule{
// 					Approval: rule.Approval{
// 						Groups: []string{"group1"},
// 						Users:  []string{"a"},
// 					},
// 				},
// 				CanRequest: true,
// 			},
// 			wantCode:                     http.StatusOK,
// 			withRequestArgumentsResponse: make(map[string]types.RequestArgument),
// 			want:                         `{"description":"","id":"","isCurrent":false,"name":"","target":{"arguments":{},"provider":{"id":"","type":""}},"timeConstraints":{"maxDurationSeconds":0},"version":""}`,
// 		},
// 		{
// 			name:           "no rule found",
// 			giveRuleID:     "notexist",
// 			mockGetRuleErr: ddb.ErrNoItems,
// 			wantCode:       http.StatusNotFound,
// 			want:           `{"error":"this rule doesn't exist or you don't have permission to access it"}`,
// 		},
// 		{
// 			name:           "not authorised to access the rule",
// 			giveRuleID:     "exists",
// 			mockGetRuleErr: rulesvc.ErrUserNotAuthorized,
// 			wantCode:       http.StatusNotFound,
// 			want:           `{"error":"this rule doesn't exist or you don't have permission to access it"}`,
// 		},
// 		{
// 			name:           "internal error",
// 			giveRuleID:     "exists",
// 			mockGetRuleErr: errors.New("internal error"),
// 			wantCode:       http.StatusInternalServerError,
// 			want:           `{"error":"Internal Server Error"}`,
// 		},
// 	}

// 	for _, tc := range testcases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			t.Parallel()
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			m := mocks.NewMockAccessRuleService(ctrl)
// 			m.EXPECT().GetRule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(tc.mockGetRuleResponse, tc.mockGetRuleErr)
// 			db := ddbmock.New(t)
// 			db.MockQuery(&storage.GetAccessRuleCurrent{Result: tc.mockGetAccessRuleVersion})
// 			db.MockQuery(&storage.ListCachedProviderOptions{Result: []cache.ProviderOption{}})
// 			if tc.withRequestArgumentsResponse != nil {
// 				m.EXPECT().RequestArguments(gomock.Any(), gomock.Any()).Return(tc.withRequestArgumentsResponse, nil)
// 			}
// 			a := API{Rules: m, DB: db}
// 			handler := newTestServer(t, &a)

// 			req, err := http.NewRequest("GET", "/api/v1/access-rules/"+tc.giveRuleID, nil)
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			req.Header.Add("Content-Type", "application/json")
// 			rr := httptest.NewRecorder()

// 			handler.ServeHTTP(rr, req)

// 			assert.Equal(t, tc.wantCode, rr.Code)

// 			data, err := io.ReadAll(rr.Body)
// 			if err != nil {
// 				t.Fatal(err)
// 			}

// 			assert.Equal(t, tc.want, string(data))
// 		})
// 	}
// }
// func TestUserGetAccessRule(t *testing.T) {
// 	type testcase struct {
// 		name                    string
// 		giveRuleID              string
// 		mockGetRuleResponse     *rule.GetAccessRuleResponse
// 		mockGetRuleErr          error
// 		mockGetGroupQueryResult *identity.Group
// 		want                    string
// 		wantCode                int
// 	}

// 	testcases := []testcase{
// 		{
// 			name:       "ok",
// 			giveRuleID: "abcd",
// 			mockGetRuleResponse: &rule.GetAccessRuleResponse{
// 				Rule: &rule.AccessRule{
// 					Approval: rule.Approval{
// 						Groups: []string{"group1"},
// 						Users:  []string{"a"},
// 					},
// 				},
// 				CanRequest: true,
// 			},
// 			mockGetGroupQueryResult: &identity.Group{
// 				ID:    "group1",
// 				Users: []string{"a", "b", "c"},
// 			},
// 			wantCode: http.StatusOK,
// 			want:     `{"next":null,"users":["a","b","c"]}`,
// 		},
// 		{
// 			name:           "no rule found",
// 			giveRuleID:     "notexist",
// 			mockGetRuleErr: ddb.ErrNoItems,
// 			wantCode:       http.StatusNotFound,
// 			want:           `{"error":"this rule doesn't exist or you don't have permission to access it"}`,
// 		},
// 		{
// 			name:           "not authorised to access the rule",
// 			giveRuleID:     "exists",
// 			mockGetRuleErr: rulesvc.ErrUserNotAuthorized,
// 			wantCode:       http.StatusNotFound,
// 			want:           `{"error":"this rule doesn't exist or you don't have permission to access it"}`,
// 		},
// 		{
// 			name:           "internal error",
// 			giveRuleID:     "exists",
// 			mockGetRuleErr: errors.New("internal error"),
// 			wantCode:       http.StatusInternalServerError,
// 			want:           `{"error":"Internal Server Error"}`,
// 		},
// 	}

// 	for _, tc := range testcases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			t.Parallel()
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			m := mocks.NewMockAccessRuleService(ctrl)
// 			m.EXPECT().GetRule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(tc.mockGetRuleResponse, tc.mockGetRuleErr)
// 			db := ddbmock.New(t)
// 			db.MockQuery(&storage.GetGroup{Result: tc.mockGetGroupQueryResult})
// 			a := API{Rules: m, DB: db}
// 			handler := newTestServer(t, &a)

// 			req, err := http.NewRequest("GET", "/api/v1/access-rules/"+tc.giveRuleID+"/approvers", nil)
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			req.Header.Add("Content-Type", "application/json")
// 			rr := httptest.NewRecorder()

// 			handler.ServeHTTP(rr, req)

// 			assert.Equal(t, tc.wantCode, rr.Code)

// 			data, err := io.ReadAll(rr.Body)
// 			if err != nil {
// 				t.Fatal(err)
// 			}

// 			assert.Equal(t, tc.want, string(data))
// 		})
// 	}
// }

// func TestLookupAccessRules(t *testing.T) {
// 	type testcase struct {
// 		name                   string
// 		giveURL                string
// 		rules                  []rule.AccessRule
// 		want                   string
// 		mockLookupRuleResponse []rulesvc.LookedUpRule
// 		mockLookupRuleErr      error
// 		wantCode               int
// 	}

// 	testcases := []testcase{
// 		{
// 			name:     "no matches",
// 			giveURL:  `/api/v1/access-rules/lookup?accountId=123456789012&permissionSetArn.label=GrantedAdministratorAccess&type=commonfate%2Faws-sso`,
// 			wantCode: http.StatusOK,
// 			rules:    nil,
// 			want:     `[]`,
// 		},
// 		{
// 			name:     "single match",
// 			giveURL:  `/api/v1/access-rules/lookup?accountId=123456789012&permissionSetArn.label=GrantedAdministratorAccess&type=commonfate%2Faws-sso`,
// 			wantCode: http.StatusOK,
// 			mockLookupRuleResponse: []rulesvc.LookedUpRule{
// 				{
// 					Rule: rule.AccessRule{
// 						ID: "test",
// 						Target: rule.Target{
// 							TargetGroupID: "test-provider",
// 							With: map[string]string{
// 								"accountId":        "123456789012",
// 								"permissionSetArn": "arn:aws:sso:::permissionSet/ssoins-1234/ps-12341",
// 							},
// 						},
// 					},
// 				},
// 			},
// 			want: `[{"accessRule":{"createdAt":"0001-01-01T00:00:00Z","description":"","id":"test","isCurrent":false,"name":"","target":{"provider":{"id":"test-provider","type":"aws-sso"}},"timeConstraints":{"maxDurationSeconds":0},"updatedAt":"0001-01-01T00:00:00Z","version":""},"selectableWithOptionValues":[{"key":"accountId","value":"123456789012"},{"key":"permissionSetArn","value":"arn:aws:sso:::permissionSet/ssoins-1234/ps-12341"}]}]`,
// 		},
// 		{
// 			name:     "multiple matches",
// 			giveURL:  `/api/v1/access-rules/lookup?accountId=123456789012&permissionSetArn.label=GrantedAdministratorAccess&type=commonfate%2Faws-sso`,
// 			wantCode: http.StatusOK,
// 			mockLookupRuleResponse: []rulesvc.LookedUpRule{
// 				{
// 					Rule: rule.AccessRule{
// 						ID: "test",
// 						Target: rule.Target{
// 							TargetGroupID: "test-provider",
// 							With: map[string]string{
// 								"accountId":        "123456789012",
// 								"permissionSetArn": "arn:aws:sso:::permissionSet/ssoins-1234/ps-12341",
// 							},
// 						},
// 					},
// 				},
// 				{
// 					Rule: rule.AccessRule{
// 						ID: "second",
// 						Target: rule.Target{
// 							TargetGroupID: "test-provider",
// 							With: map[string]string{
// 								"accountId":        "123456789012",
// 								"permissionSetArn": "arn:aws:sso:::permissionSet/ssoins-1234/ps-12341",
// 							},
// 						},
// 					},
// 				},
// 			},
// 			want: `[{"accessRule":{"createdAt":"0001-01-01T00:00:00Z","description":"","id":"test","isCurrent":false,"name":"","target":{"provider":{"id":"test-provider","type":"aws-sso"}},"timeConstraints":{"maxDurationSeconds":0},"updatedAt":"0001-01-01T00:00:00Z","version":""},"selectableWithOptionValues":[{"key":"accountId","value":"123456789012"},{"key":"permissionSetArn","value":"arn:aws:sso:::permissionSet/ssoins-1234/ps-12341"}]}]`,
// 		},
// 		{
// 			name:     "match with selectable",
// 			giveURL:  `/api/v1/access-rules/lookup?accountId=123456789012&permissionSetArn.label=GrantedAdministratorAccess&type=commonfate%2Faws-sso`,
// 			wantCode: http.StatusOK,
// 			mockLookupRuleResponse: []rulesvc.LookedUpRule{
// 				{
// 					Rule: rule.AccessRule{
// 						ID: "test",
// 						Target: rule.Target{
// 							TargetGroupID: "test-provider",
// 							// WithSelectable: map[string][]string{
// 							// 	"accountId":        {"123456789012", "other"},
// 							// 	"permissionSetArn": {"arn:aws:sso:::permissionSet/ssoins-1234/ps-12341", "other"},
// 							// },
// 						},
// 					},
// 					SelectableWithOptionValues: []types.KeyValue{
// 						{
// 							Key:   "accountId",
// 							Value: "123456789012",
// 						},
// 						{
// 							Key:   "permissionSetArn",
// 							Value: "arn:aws:sso:::permissionSet/ssoins-1234/ps-12341",
// 						},
// 					},
// 				},
// 			},
// 			want: `[{"accessRule":{"createdAt":"0001-01-01T00:00:00Z","description":"","id":"test","isCurrent":false,"name":"","target":{"provider":{"id":"test-provider","type":"aws-sso"}},"timeConstraints":{"maxDurationSeconds":0},"updatedAt":"0001-01-01T00:00:00Z","version":""},"selectableWithOptionValues":[{"key":"accountId","value":"123456789012"},{"key":"permissionSetArn","value":"arn:aws:sso:::permissionSet/ssoins-1234/ps-12341"}]}]`,
// 		},
// 	}

// 	for _, tc := range testcases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			t.Parallel()
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			m := mocks.NewMockAccessRuleService(ctrl)
// 			// m.EXPECT().LookupRule(gomock.Any(), gomock.Any()).Return(tc.mockLookupRuleResponse, tc.mockLookupRuleErr)

// 			a := API{Rules: m}
// 			handler := newTestServer(t, &a)

// 			req, err := http.NewRequest("GET", tc.giveURL, nil)
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			req.Header.Add("Content-Type", "application/json")
// 			rr := httptest.NewRecorder()

// 			handler.ServeHTTP(rr, req)

// 			assert.Equal(t, tc.wantCode, rr.Code)

// 			data, err := io.ReadAll(rr.Body)
// 			if err != nil {
// 				t.Fatal(err)
// 			}

// 			assert.Equal(t, tc.want, string(data))
// 		})
// 	}
// }
