package rulesvc

import (
	"context"
	"net/http"
	"testing"

	"github.com/benbjohnson/clock"
	ahTypes "github.com/common-fate/common-fate/accesshandler/pkg/types"
	"github.com/common-fate/common-fate/accesshandler/pkg/types/ahmocks"
	"github.com/common-fate/common-fate/pkg/cache"
	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/common-fate/common-fate/pkg/service/rulesvc/mocks"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb/ddbmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUpdateAccessRule(t *testing.T) {

	type testcase struct {
		name            string
		givenUserID     string
		givenRule       rule.AccessRule
		givenUpdateBody types.CreateAccessRuleRequest
		wantErr         error
		want            *rule.AccessRule
	}

	in := types.CreateAccessRuleRequest{}

	ruleID := "override"
	versionID := types.NewVersionID()
	userID := "user1"
	clk := clock.NewMock()
	now := clk.Now()
	cacheArgOptionsResponse := []cache.ProviderOption{}
	cacheArgGroupOptionsResponse := []cache.ProviderArgGroupOption{}

	/**
	Input values needed:
	- UpdateOpts.Rule
	- UpdateOpts.UpdateRequest
	*/
	mockRule := rule.AccessRule{
		ID:       ruleID,
		Approval: rule.Approval(in.Approval),
		Status:   rule.ACTIVE,
		Metadata: rule.AccessRuleMetadata{
			CreatedAt: now,
			CreatedBy: userID,
			UpdatedAt: now,
			UpdatedBy: userID,
		},
		Target: rule.Target{
			ProviderID:     "hello",
			ProviderType:   "awssso",
			With:           map[string]string{},
			WithSelectable: map[string][]string{},
		},
	}

	mockRuleUpdateBody := types.CreateAccessRuleRequest{
		Approval: types.ApproverConfig{
			Users: []string{"user1", "user2"},
		},
		Name:        "changing the name",
		Description: "changing the description name",
		Groups:      []string{"group1", "group2"},
		TimeConstraints: types.TimeConstraints{
			MaxDurationSeconds: 600,
		},
		Target: types.CreateAccessRuleTarget{
			ProviderId: "newTarget",
			With:       make(map[string]types.CreateAccessRuleTargetDetailArguments),
		},
	}

	want := rule.AccessRule{
		ID: ruleID,
		Approval: rule.Approval{
			Users: mockRuleUpdateBody.Approval.Users,
		},
		Status:      rule.ACTIVE,
		Description: mockRuleUpdateBody.Description,
		Name:        mockRuleUpdateBody.Name,
		Groups:      mockRuleUpdateBody.Groups,
		Metadata: rule.AccessRuleMetadata{
			CreatedAt: now,
			CreatedBy: userID,
			UpdatedAt: now,
			UpdatedBy: userID,
		},
		TimeConstraints: types.TimeConstraints{
			MaxDurationSeconds: 600,
		},
		Version: versionID,
		Target: rule.Target{
			ProviderID:               "newTarget",
			ProviderType:             "awssso",
			With:                     make(map[string]string),
			WithSelectable:           make(map[string][]string),
			WithArgumentGroupOptions: make(map[string]map[string][]string),
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
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {

			dbc := ddbmock.Client{
				PutBatchErr: tc.wantErr,
			}
			clk := clock.NewMock()
			ctrl := gomock.NewController(t)

			defer ctrl.Finish()

			m := ahmocks.NewMockClientWithResponsesInterface(ctrl)
			m.EXPECT().GetProviderWithResponse(gomock.Any(), gomock.Eq(tc.givenUpdateBody.Target.ProviderId)).Return(&ahTypes.GetProviderResponse{HTTPResponse: &http.Response{StatusCode: 200}, JSON200: &ahTypes.Provider{Type: "awssso"}}, nil)
			m.EXPECT().GetProviderArgsWithResponse(gomock.Any(), gomock.Eq(tc.givenUpdateBody.Target.ProviderId)).Return(&ahTypes.GetProviderArgsResponse{HTTPResponse: &http.Response{StatusCode: 200}, JSON200: &ahTypes.ArgSchema{}}, nil)

			cm := mocks.NewMockCacheService(ctrl)
			cm.EXPECT().LoadCachedProviderArgOptions(gomock.Any(), gomock.Eq(tc.givenUpdateBody.Target.ProviderId), gomock.Any()).AnyTimes().Return(false, cacheArgOptionsResponse, cacheArgGroupOptionsResponse, nil)

			s := Service{
				Clock:    clk,
				DB:       &dbc,
				AHClient: m,
				Cache:    cm,
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
				assert.NotEmpty(t, got.Version)
				assert.NotEmpty(t, got.ID)
				got.ID = ruleID
				got.Version = versionID
			}

			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.want, got)

		})
	}

}
