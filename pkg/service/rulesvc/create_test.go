package rulesvc

import (
	"context"
	"errors"
	"testing"

	"github.com/benbjohnson/clock"
	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/common-fate/common-fate/pkg/service/rulesvc/mocks"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/target"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
	"github.com/common-fate/ddb/ddbmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateAccessRule(t *testing.T) {
	type testcase struct {
		name               string
		givenUserID        string
		give               types.CreateAccessRuleRequest
		wantErr            error
		want               *rule.AccessRule
		wantTargetGroup    target.Group
		wantTargetGroupErr error
	}

	in := types.CreateAccessRuleRequest{
		Approval: types.AccessRuleApproverConfig{
			Groups: []string{"test"},
			Users:  []string{"test"},
		},
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
					Schema:    target.GroupSchema{},
					Icon:      "",
					CreatedAt: now,
					UpdatedAt: now,
				},
				FieldFilterExpessions: map[string]rule.FieldFilterExpessions{},
			},
		},

		TimeConstraints: in.TimeConstraints,
	}

	mockRuleLongerThan6months := in
	mockRuleLongerThan6months.TimeConstraints = types.AccessRuleTimeConstraints{MaxDurationSeconds: 26*7*24*3600 + 1}

	/**
	There are two test cases here:
	- Create a valid rule
	*/
	testcases := []testcase{
		{
			name:        "ok",
			givenUserID: userID,
			give:        in,
			want:        &mockRule,
			wantTargetGroup: target.Group{
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
		},
		{
			name:        "max duration longer than 6 months",
			givenUserID: userID,
			give:        mockRuleLongerThan6months,
			wantErr:     errors.New("access rule cannot be longer than 6 months"),
			wantTargetGroup: target.Group{
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
		},
		{
			name:               "target group not found errors gracefully",
			givenUserID:        userID,
			give:               in,
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

			mockCache := mocks.NewMockCacheService(ctrl)
			if tc.wantTargetGroupErr == nil && tc.wantErr == nil {
				mockCache.EXPECT().RefreshCachedTargets(gomock.Any()).Return(nil)

			}

			s := Service{
				Clock: clk,
				DB:    dbc,
				Cache: mockCache,
			}

			got, err := s.CreateAccessRule(context.Background(), tc.givenUserID, tc.give)

			// This is the only thing from service layer that we can't mock yet, hence the override
			if err == nil {
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

func TestProcessTarget(t *testing.T) {
	type testcase struct {
		name                     string
		give                     []types.CreateAccessRuleTarget
		wantErr                  error
		want                     []rule.Target
		wantTargetGroupLookup    *target.Group
		wantTargetGroupLookupErr error
	}

	tg1 := target.Group{
		ID: "123",
		From: target.From{
			Publisher: "test",
			Name:      "test",
			Version:   "test",
			Kind:      "test",
		},
		Schema: providerregistrysdk.Target{},
	}

	testcases := []testcase{
		{
			name: "ok ",
			give: []types.CreateAccessRuleTarget{
				{
					TargetGroupId:         "123",
					FieldFilterExpessions: make(map[string]interface{}),
				},
			},
			wantTargetGroupLookup: &tg1,
			want: []rule.Target{
				{
					TargetGroup:           tg1,
					FieldFilterExpessions: map[string]rule.FieldFilterExpessions{},
				},
			},
		},
		{
			name: "duplicate target group fails ",
			give: []types.CreateAccessRuleTarget{
				{
					TargetGroupId:         "123",
					FieldFilterExpessions: make(map[string]interface{}),
				},
				{
					TargetGroupId:         "123",
					FieldFilterExpessions: make(map[string]interface{}),
				},
			},
			wantTargetGroupLookup: &tg1,
			want:                  nil,
			wantErr:               errors.New("duplicate target in access rule"),
		},
		{
			name: "target group not found",
			give: []types.CreateAccessRuleTarget{
				{
					TargetGroupId:         "123",
					FieldFilterExpessions: make(map[string]interface{}),
				},
			},
			wantTargetGroupLookup:    nil,
			want:                     nil,
			wantTargetGroupLookupErr: ddb.ErrNoItems,
			wantErr:                  ddb.ErrNoItems,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			dbc := ddbmock.New(t)

			dbc.MockQueryWithErr(&storage.GetTargetGroup{Result: tc.wantTargetGroupLookup}, tc.wantTargetGroupLookupErr)

			clk := clock.NewMock()
			ctrl := gomock.NewController(t)

			defer ctrl.Finish()

			s := Service{
				Clock: clk,
				DB:    dbc,
			}
			got, err := s.ProcessTargets(context.Background(), tc.give)
			if tc.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.wantErr.Error())
			}
			assert.Equal(t, tc.want, got)
		})
	}
}
