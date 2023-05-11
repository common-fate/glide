package api

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/benbjohnson/clock"
	"github.com/common-fate/common-fate/pkg/api/mocks"
	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/common-fate/common-fate/pkg/service/rulesvc"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/target"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
	"github.com/common-fate/ddb/ddbmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAdminCreateAccessRule(t *testing.T) {
	type testcase struct {
		name          string
		give          string
		mockCreate    *rule.AccessRule
		mockCreateErr error

		//idpUser  *types.User
		wantCode int
		wantBody string
	}

	testcases := []testcase{
		{
			name: "ok",
			give: `{"priority":4,"approval":{"groups":["group1","group2"],"users":["user1","user2"]},"description":"Test Access Rule","groups":["group_a","group_b"],"name":"Test Access Rule","targets":[{"fieldFilterExpessions":{"field1":"value1","field2":"value2"},"targetGroupId":"target_group_id"}],"timeConstraints":{"maxDurationSeconds":3600}}`,
			mockCreate: &rule.AccessRule{
				ID: "rule1",

				Description: "string",
				Name:        "string",
				Groups:      []string{"string"},
				Targets: []rule.Target{
					{
						TargetGroup: target.Group{
							ID: "123",
							From: target.From{
								Name:      "test",
								Publisher: "commonfate",
								Version:   "v1.1.1",
								Kind:      "Account",
							},
						},
						FieldFilterExpessions: map[string][]types.Operation{},
					},
				},
				Priority: 4,
			},
			wantCode: http.StatusCreated,
			// idpUser: &types.User{
			// 	Id:    "userid",
			// 	Email: "test@test.com",
			// },

			wantBody: `{"approval":{"groups":[],"users":[]},"description":"string","groups":["string"],"id":"rule1","metadata":{"createdAt":"0001-01-01T00:00:00Z","createdBy":"","updatedAt":"0001-01-01T00:00:00Z","updatedBy":""},"name":"string","priority":4,"targets":[{"fieldFilterExpessions":{},"targetGroup":{"createdAt":"0001-01-01T00:00:00Z","from":{"kind":"Account","name":"test","publisher":"commonfate","version":"v1.1.1"},"icon":"","id":"123","schema":{},"updatedAt":"0001-01-01T00:00:00Z"}}],"timeConstraints":{"maxDurationSeconds":0}}`,
		},
		{
			name:          "id already exists",
			give:          `{"priority":4,"approval":{"groups":["group1","group2"],"users":["user1","user2"]},"description":"Test Access Rule","groups":["group_a","group_b"],"name":"Test Access Rule","targets":[{"fieldFilterExpessions":{"field1":"value1","field2":"value2"},"targetGroupId":"target_group_id"}],"timeConstraints":{"maxDurationSeconds":3600}}`,
			mockCreateErr: rulesvc.ErrRuleIdAlreadyExists,
			wantCode:      http.StatusBadRequest,
			wantBody:      `{"error":"access rule id already exists"}`,
		},
		{
			name:     "fail when rule doesn't meet maxduration req",
			give:     `{"priority":4,"approval":{"groups":["group1","group2"],"users":["user1","user2"]},"description":"Test Access Rule","groups":["group_a","group_b"],"name":"Test Access Rule","targets":[{"fieldFilterExpessions":{"field1":"value1","field2":"value2"},"targetGroupId":"target_group_id"}],"timeConstraints":{"maxDurationSeconds":6}}`,
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

			a := API{Rules: m}
			handler := newTestServer(t, &a)

			req, err := http.NewRequest("POST", "/api/v1/admin/access-rules", strings.NewReader(tc.give))
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

func TestAdminUpdateAccessRule(t *testing.T) {
	type testcase struct {
		name          string
		give          string
		mockCreate    *rule.AccessRule
		mockCreateErr error

		//idpUser  *types.User
		wantCode int
		wantBody string
		wantErr  string
	}

	testcases := []testcase{
		{
			name: "ok",
			give: `{"priority":4,"approval":{"groups":["group1","group2"],"users":["user1","user2"]},"description":"Test Access Rule","groups":["group_a","group_b"],"name":"Test Access Rule","targets":[{"fieldFilterExpessions":{"field1":"value1","field2":"value2"},"targetGroupId":"target_group_id"}],"timeConstraints":{"maxDurationSeconds":3600}}`,
			mockCreate: &rule.AccessRule{
				ID:          "rule1",
				Description: "string",
				Name:        "string",
				Groups:      []string{"string"},
				Targets: []rule.Target{
					{
						TargetGroup: target.Group{
							ID: "123",
							From: target.From{
								Name:      "test",
								Publisher: "commonfate",
								Version:   "v1.1.1",
								Kind:      "Account",
							},
						},
						FieldFilterExpessions: map[string][]types.Operation{},
					},
				},
				Priority: 4,
			},
			wantCode: http.StatusAccepted,
			wantBody: `{"approval":{"groups":[],"users":[]},"description":"string","groups":["string"],"id":"rule1","metadata":{"createdAt":"0001-01-01T00:00:00Z","createdBy":"","updatedAt":"0001-01-01T00:00:00Z","updatedBy":""},"name":"string","priority":4,"targets":[{"fieldFilterExpessions":{},"targetGroup":{"createdAt":"0001-01-01T00:00:00Z","from":{"kind":"Account","name":"test","publisher":"commonfate","version":"v1.1.1"},"icon":"","id":"123","schema":{},"updatedAt":"0001-01-01T00:00:00Z"}}],"timeConstraints":{"maxDurationSeconds":0}}`,
		},

		{
			name:     "malformed",
			give:     `malformed json input`,
			wantCode: http.StatusBadRequest,
			wantErr:  `{"error":"request body has an error: failed to decode request body: invalid character 'm' looking for beginning of value"}`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mocks.NewMockAccessRuleService(ctrl)
			if tc.mockCreate != nil {
				m.EXPECT().UpdateRule(gomock.Any(), gomock.Any()).Return(tc.mockCreate, tc.mockCreateErr)
			}
			db := ddbmock.New(t)
			db.MockQuery(&storage.GetAccessRule{Result: tc.mockCreate})
			a := API{Rules: m, DB: db}
			handler := newTestServer(t, &a)

			req, err := http.NewRequest("PUT", "/api/v1/admin/access-rules/"+"rule1", strings.NewReader(tc.give))
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
			if tc.wantErr != "" {
				assert.Equal(t, tc.wantErr, string(data))
				return
			}
			if tc.wantBody != "" {
				assert.Equal(t, tc.wantBody, string(data))
			}
		})
	}
}

func TestAdminListAccessRules(t *testing.T) {
	type testcase struct {
		name string

		rules       []rule.AccessRule
		want        string
		mockListErr error
		wantCode    int
	}
	clk := clock.NewMock()
	now := clk.Now()

	testcases := []testcase{
		{
			name: "ok",
			//mockListErr: types.ErrNoGroupsFound,
			wantCode: http.StatusOK,
			rules: []rule.AccessRule{
				{
					ID:          "rule1",
					Description: "string",
					Name:        "string",
					Groups:      []string{"string"},

					// This should not be included in the response for users
					Approval: rule.Approval{
						Groups: []string{"a"},
						Users:  []string{"b"},
					},
					Targets: []rule.Target{
						{
							TargetGroup: target.Group{
								ID: "123",
								From: target.From{
									Name:      "test",
									Publisher: "commonfate",
									Kind:      "Account",
									Version:   "v1.1.1",
								},
								Schema:    target.GroupSchema{},
								Icon:      "",
								CreatedAt: now,
								UpdatedAt: now,
							},
							FieldFilterExpessions: map[string][]types.Operation{},
						},
					},
				},
				{
					ID: "rule2",

					Description: "string",
					Name:        "string",
					Groups:      []string{"string"},

					// This should not be included in the response for users
					Approval: rule.Approval{
						Groups: []string{"a"},
						Users:  []string{"b"},
					},
					Targets: []rule.Target{
						{
							TargetGroup: target.Group{
								ID: "123",
								From: target.From{
									Name:      "test",
									Publisher: "commonfate",
									Kind:      "Account",
									Version:   "v1.1.1",
								},
								Schema:    target.GroupSchema{},
								Icon:      "",
								CreatedAt: now,
								UpdatedAt: now,
							},
							FieldFilterExpessions: map[string][]types.Operation{},
						},
					},
				},
			},

			want: `{"accessRules":[{"approval":{"groups":["a"],"users":["b"]},"description":"string","groups":["string"],"id":"rule1","metadata":{"createdAt":"0001-01-01T00:00:00Z","createdBy":"","updatedAt":"0001-01-01T00:00:00Z","updatedBy":""},"name":"string","status":"ACTIVE","targets":[{"fieldFilterExpessions":{},"targetGroup":{"createdAt":"1970-01-01T10:00:00+10:00","from":{"kind":"Account","name":"test","publisher":"commonfate","version":"v1.1.1"},"icon":"","id":"123","schema":{},"updatedAt":"1970-01-01T10:00:00+10:00"}}],"timeConstraints":{"maxDurationSeconds":0}},{"approval":{"groups":["a"],"users":["b"]},"description":"string","groups":["string"],"id":"rule2","metadata":{"createdAt":"0001-01-01T00:00:00Z","createdBy":"","updatedAt":"0001-01-01T00:00:00Z","updatedBy":""},"name":"string","status":"ACTIVE","targets":[{"fieldFilterExpessions":{},"targetGroup":{"createdAt":"1970-01-01T10:00:00+10:00","from":{"kind":"Account","name":"test","publisher":"commonfate","version":"v1.1.1"},"icon":"","id":"123","schema":{},"updatedAt":"1970-01-01T10:00:00+10:00"}}],"timeConstraints":{"maxDurationSeconds":0}}],"next":null}`,
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
			db.MockQueryWithErr(&storage.ListAccessRulesByPriority{Result: tc.rules}, tc.mockListErr)

			a := API{DB: db}
			handler := newTestServer(t, &a)

			req, err := http.NewRequest("GET", "/api/v1/admin/access-rules", nil)
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
