package rulesvc

import (
	"context"
	"testing"

	"github.com/benbjohnson/clock"
	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/target"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
	"github.com/common-fate/ddb/ddbmock"
	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUpdateAccessRule(t *testing.T) {

	type testcase struct {
		name               string
		givenUserID        string
		givenRule          rule.AccessRule
		givenUpdateBody    types.CreateAccessRuleRequest
		wantErr            error
		want               *rule.AccessRule
		wantTargetGroup    target.Group
		wantTargetGroupErr error
	}

	in := types.CreateAccessRuleRequest{
		Approval:        types.AccessRuleApproverConfig{},
		Description:     "test",
		Name:            "test",
		Groups:          []string{"group_a"},
		TimeConstraints: types.AccessRuleTimeConstraints{MaxDurationSeconds: 3600},

		Targets: []types.CreateAccessRuleTarget{
			{
				TargetGroupId:         "test",
				FieldFilterExpessions: make(map[string]interface{}),
			},
		},
	}

	ruleID := "override"
	userID := "user1"
	clk := clock.NewMock()
	now := clk.Now()

	mockRule := rule.AccessRule{
		ID:          ruleID,
		Approval:    rule.Approval(in.Approval),
		Description: in.Description,
		Name:        in.Name,
		Groups:      in.Groups,
		Metadata: rule.AccessRuleMetadata{
			CreatedAt: now,
			CreatedBy: userID,
			UpdatedAt: now,
			UpdatedBy: userID,
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
					Schema:    providerregistrysdk.Target{},
					Icon:      "",
					CreatedAt: now,
					UpdatedAt: now,
				},
				FieldFilterExpessions: map[string]rule.FieldFilterExpessions{},
			},
		},

		TimeConstraints: in.TimeConstraints,
	}

	mockRuleUpdateBody := types.CreateAccessRuleRequest{
		Approval:        types.AccessRuleApproverConfig{},
		Description:     "updated description",
		Name:            "updated name",
		Groups:          []string{"group_b"},
		TimeConstraints: types.AccessRuleTimeConstraints{MaxDurationSeconds: 3601},

		Targets: []types.CreateAccessRuleTarget{
			{
				TargetGroupId:         "test",
				FieldFilterExpessions: make(map[string]interface{}),
			},
		},
	}

	want := rule.AccessRule{
		ID: ruleID,
		Approval: rule.Approval{
			Users: mockRuleUpdateBody.Approval.Users,
		},
		Description: mockRuleUpdateBody.Description,
		Name:        mockRuleUpdateBody.Name,
		Groups:      mockRuleUpdateBody.Groups,
		Metadata: rule.AccessRuleMetadata{
			CreatedAt: now,
			CreatedBy: userID,
			UpdatedAt: now,
			UpdatedBy: userID,
		},
		TimeConstraints: types.AccessRuleTimeConstraints{
			MaxDurationSeconds: 3601,
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
					Schema:    providerregistrysdk.Target{},
					Icon:      "",
					CreatedAt: now,
					UpdatedAt: now,
				},
				FieldFilterExpessions: map[string]rule.FieldFilterExpessions{},
			},
		},
	}

	/**
	Things to test:
	- Control test case (pass) ✅
	- Non admin user cannot update rule ✅
	*/
	testcases := []testcase{
		{
			name:            "ok",
			givenUserID:     userID,
			givenRule:       mockRule,
			givenUpdateBody: mockRuleUpdateBody,
			want:            &want,
			wantTargetGroup: target.Group{
				ID: "123",
				From: target.From{
					Name:      "test",
					Publisher: "commonfate",
					Kind:      "Account",
					Version:   "v1.1.1",
				},
				Schema:    providerregistrysdk.Target{},
				Icon:      "",
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		{
			name:               "target group not found fails gracefully",
			givenUserID:        userID,
			givenRule:          mockRule,
			givenUpdateBody:    mockRuleUpdateBody,
			wantTargetGroupErr: ddb.ErrNoItems,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {

			dbc := ddbmock.New(t)

			clk := clock.NewMock()
			ctrl := gomock.NewController(t)

			defer ctrl.Finish()

			dbc.MockQueryWithErr(&storage.GetTargetGroup{Result: &tc.wantTargetGroup}, tc.wantTargetGroupErr)

			s := Service{
				Clock: clk,
				DB:    dbc,
			}

			got, err := s.UpdateRule(context.Background(), &UpdateOpts{
				UpdaterID:      tc.givenUserID,
				Rule:           tc.givenRule,
				UpdateRequest:  tc.givenUpdateBody,
				ApprovalGroups: []rule.Approval{},
			})

			// This is the only thing from service layer that we can't mock yet, hence the override
			if err == nil {
				// Rule id and version id must not be empty strings, we check this prior to overwriting them
				assert.NotEmpty(t, got.ID)
				got.ID = ruleID
			}

			if tc.wantTargetGroupErr != nil {
				assert.Equal(t, tc.wantTargetGroupErr.Error(), err.Error())
				return
			}

			if err != nil {
				assert.Equal(t, tc.wantErr.Error(), err.Error())
			}

			assert.Equal(t, tc.want, got)

		})
	}

}
