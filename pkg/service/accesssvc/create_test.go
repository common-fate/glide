package accesssvc

import (
	"context"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/cache"
	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/common-fate/common-fate/pkg/service/accesssvc/mocks"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb/ddbmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateRequest(t *testing.T) {
	type testcase struct {
		name                   string
		user                   identity.User
		createRequest          types.CreateAccessRequestRequest
		withMockPreflight      *access.Preflight
		withMockPreflightErr   error
		withMockGetAccessRules []rule.AccessRule
		withMockGetApprovers   [][]string
		want                   *access.RequestWithGroupsWithTargets
		wantErr                error
	}

	reason := "test_reason"

	clk := clock.NewMock()
	user := identity.User{
		ID:        types.NewUserID(),
		FirstName: "tester",
		LastName:  "wow",
		Email:     "test@example.com",
	}
	requestedBy := access.RequestedBy{
		ID:        user.ID,
		FirstName: "tester",
		LastName:  "wow",
		Email:     "test@example.com",
	}
	testcases := []testcase{
		{
			name: "ok",
			user: user,
			createRequest: types.CreateAccessRequestRequest{
				GroupOptions: []types.CreateAccessRequestGroupOptions{
					{
						Id: "group",
						Timing: types.RequestAccessGroupTiming{
							DurationSeconds: 3600,
						},
					},
				},
				Reason: &reason,
			},
			withMockPreflight: &access.Preflight{
				AccessGroups: []access.PreflightAccessGroup{
					{
						ID: "group",
						Targets: []access.PreflightAccessGroupTarget{
							{
								Target: cache.Target{
									Kind: cache.Kind{
										Publisher: "publisher",
										Name:      "name",
										Kind:      "kind",
										Icon:      "icon",
									},
									Fields: []cache.Field{{ID: "a"}},
								},
							},
						},
					},
				},
			},
			withMockGetAccessRules: []rule.AccessRule{
				{},
			},
			withMockGetApprovers: [][]string{{}},
			want: &access.RequestWithGroupsWithTargets{
				Request: access.Request{
					RequestedBy:      requestedBy,
					CreatedAt:        clk.Now(),
					RequestStatus:    types.PENDING,
					GroupTargetCount: 1,
					Purpose:          access.Purpose{Reason: &reason},
				},
				Groups: []access.GroupWithTargets{
					{
						Group: access.Group{
							CreatedAt:     clk.Now(),
							UpdatedAt:     clk.Now(),
							RequestStatus: types.PENDING,
							Status:        types.RequestAccessGroupStatusPENDINGAPPROVAL,
							RequestedBy:   requestedBy,
							RequestedTiming: access.Timing{
								Duration: 3600 * time.Second,
							},
							RequestPurposeReason: reason,
						},
						Targets: []access.GroupTarget{
							{
								CreatedAt:     clk.Now(),
								UpdatedAt:     clk.Now(),
								RequestStatus: types.PENDING,
								RequestedBy:   requestedBy,
								TargetKind: cache.Kind{
									Publisher: "publisher",
									Name:      "name",
									Kind:      "kind",
									Icon:      "icon",
								},
								TargetCacheID: "publisher#name#kind#a##",
								Fields:        []access.Field{{ID: "a", Value: access.FieldValue{Type: "string"}}},
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			db := ddbmock.New(t)
			db.MockQueryWithErr(&storage.GetPreflight{Result: tc.withMockPreflight}, tc.withMockPreflightErr)
			for i := range tc.withMockGetAccessRules {
				db.MockQuery(&storage.GetAccessRule{Result: &tc.withMockGetAccessRules[i]})
			}

			ctrl := gomock.NewController(t)
			ep := mocks.NewMockEventPutter(ctrl)
			ep.EXPECT().Put(gomock.Any(), gomock.Any()).AnyTimes()

			rs := mocks.NewMockAccessRuleService(ctrl)
			for _, ap := range tc.withMockGetApprovers {
				rs.EXPECT().GetApprovers(gomock.Any(), gomock.Any()).Return(ap, nil)
			}

			s := Service{
				Clock:       clk,
				DB:          db,
				EventPutter: ep,
				Rules:       rs,
			}
			got, err := s.CreateRequest(context.Background(), tc.user, tc.createRequest)
			if tc.wantErr != nil {
				assert.EqualError(t, err, tc.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}

			// Overwrite all the IDs
			got.Request.ID = ""
			for i, g := range got.Groups {
				g.Group.ID = ""
				g.Group.RequestID = ""
				for i, t := range g.Targets {
					t.ID = ""
					t.RequestID = ""
					t.GroupID = ""
					g.Targets[i] = t
				}
				got.Groups[i] = g
			}
			assert.Equal(t, tc.want, got)
		})
	}

}
