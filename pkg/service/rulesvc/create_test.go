package rulesvc

import (
	"context"
	"errors"
	"testing"

	"github.com/benbjohnson/clock"
	"github.com/common-fate/common-fate/pkg/rule"
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

			s := Service{
				Clock: clk,
				DB:    dbc,
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

// func TestProcessTarget(t *testing.T) {
// 	type testcase struct {
// 		name                     string
// 		give                     types.CreateAccessRuleTarget
// 		wantErr                  error
// 		withProviderResponse     types.GetProviderResponse
// 		withProviderArgsResponse types.GetProviderArgsResponse
// 		dontExpectCacheCall      bool
// 		want                     rule.Target
// 	}

// 	cacheArgOptionsResponse := []cache.ProviderOption{{Provider: "abcd", Arg: "accountId", Label: "", Value: "account1"}, {Provider: "abcd", Arg: "accountId", Label: "", Value: "account2"}, {Provider: "abcd", Arg: "permissionSetArn", Label: "", Value: "abcdefg"}}
// 	cacheArgGroupOptionsResponse := []cache.ProviderArgGroupOption{{Provider: "abcd", Arg: "accountId", Group: "organizationalUnit", Label: "", Value: "orgunit1"}, {Provider: "abcd", Arg: "accountId", Group: "organizationalUnit", Label: "", Value: "orgunit2"}}
// 	ssov2Schema := (&ssov2.Provider{}).ArgSchema().ToAPI()
// 	testVaultSchema := (&testvault.Provider{}).ArgSchema().ToAPI()

// 	testcases := []testcase{
// 		{
// 			name: "ok testvault with input element",
// 			give: types.CreateAccessRuleTarget{
// 				ProviderId: "abcd",
// 				With: types.CreateAccessRuleTarget_With{
// 					AdditionalProperties: map[string]types.CreateAccessRuleTargetDetailArguments{
// 						"vault": {
// 							Groupings: types.CreateAccessRuleTargetDetailArguments_Groupings{
// 								AdditionalProperties: map[string][]string{},
// 							},
// 							Values: []string{"example-vault"},
// 						},
// 					},
// 				},
// 			},
// 			want: rule.Target{
// 				ProviderID:               "abcd",
// 				BuiltInProviderType:      "testvault",
// 				With:                     map[string]string{"vault": "example-vault"},
// 				WithSelectable:           map[string][]string{},
// 				WithArgumentGroupOptions: map[string]map[string][]string{},
// 			},
// 			withProviderResponse: types.GetProviderResponse{
// 				HTTPResponse: &http.Response{
// 					StatusCode: http.StatusOK,
// 				},
// 				JSON200: &types.Provider{
// 					Id:   "abcd",
// 					Type: "testvault",
// 				},
// 			},
// 			withProviderArgsResponse: types.GetProviderArgsResponse{
// 				HTTPResponse: &http.Response{
// 					StatusCode: http.StatusOK,
// 				},
// 				JSON200: &testVaultSchema,
// 			},
// 			dontExpectCacheCall: true,
// 		},
// 		{
// 			name: "ok single value for field is stored in target.With, all other fields on target are empty",
// 			give: types.CreateAccessRuleTarget{
// 				ProviderId: "abcd",
// 				With: types.CreateAccessRuleTarget_With{
// 					AdditionalProperties: map[string]types.CreateAccessRuleTargetDetailArguments{
// 						"accountId": {
// 							Groupings: types.CreateAccessRuleTargetDetailArguments_Groupings{
// 								AdditionalProperties: map[string][]string{},
// 							},
// 							Values: []string{"account1"},
// 						},
// 						"permissionSetArn": {
// 							Groupings: types.CreateAccessRuleTargetDetailArguments_Groupings{
// 								AdditionalProperties: map[string][]string{},
// 							},
// 							Values: []string{"abcdefg"},
// 						},
// 					},
// 				},
// 			},
// 			want: rule.Target{
// 				ProviderID:               "abcd",
// 				BuiltInProviderType:      "awssso",
// 				With:                     map[string]string{"accountId": "account1", "permissionSetArn": "abcdefg"},
// 				WithSelectable:           map[string][]string{},
// 				WithArgumentGroupOptions: map[string]map[string][]string{},
// 			},
// 			withProviderResponse: types.GetProviderResponse{
// 				HTTPResponse: &http.Response{
// 					StatusCode: http.StatusOK,
// 				},
// 				JSON200: &types.Provider{
// 					Id:   "abcd",
// 					Type: "awssso",
// 				},
// 			},
// 			withProviderArgsResponse: types.GetProviderArgsResponse{
// 				HTTPResponse: &http.Response{
// 					StatusCode: http.StatusOK,
// 				},
// 				JSON200: &ssov2Schema,
// 			},
// 		},
// 		{
// 			name: "ok single value for field is stored in target.WithSelectable, all other fields on target are empty",
// 			give: types.CreateAccessRuleTarget{
// 				ProviderId: "abcd",
// 				With: types.CreateAccessRuleTarget_With{
// 					AdditionalProperties: map[string]types.CreateAccessRuleTargetDetailArguments{
// 						"accountId": {
// 							Groupings: types.CreateAccessRuleTargetDetailArguments_Groupings{
// 								AdditionalProperties: map[string][]string{},
// 							},
// 							Values: []string{"account1", "account2"},
// 						},
// 						"permissionSetArn": {
// 							Groupings: types.CreateAccessRuleTargetDetailArguments_Groupings{
// 								AdditionalProperties: map[string][]string{},
// 							},
// 							Values: []string{"abcdefg"},
// 						},
// 					},
// 				},
// 			},
// 			want: rule.Target{
// 				ProviderID:               "abcd",
// 				BuiltInProviderType:      "awssso",
// 				With:                     map[string]string{"permissionSetArn": "abcdefg"},
// 				WithSelectable:           map[string][]string{"accountId": {"account1", "account2"}},
// 				WithArgumentGroupOptions: map[string]map[string][]string{},
// 			},
// 			withProviderResponse: types.GetProviderResponse{
// 				HTTPResponse: &http.Response{
// 					StatusCode: http.StatusOK,
// 				},
// 				JSON200: &types.Provider{
// 					Id:   "abcd",
// 					Type: "awssso",
// 				},
// 			},
// 			withProviderArgsResponse: types.GetProviderArgsResponse{
// 				HTTPResponse: &http.Response{
// 					StatusCode: http.StatusOK,
// 				},
// 				JSON200: &ssov2Schema,
// 			},
// 		},
// 		{
// 			name: "ok group is provided for one of teh arguments",
// 			give: types.CreateAccessRuleTarget{
// 				ProviderId: "abcd",
// 				With: types.CreateAccessRuleTarget_With{
// 					AdditionalProperties: map[string]types.CreateAccessRuleTargetDetailArguments{
// 						"accountId": {
// 							Groupings: types.CreateAccessRuleTargetDetailArguments_Groupings{
// 								AdditionalProperties: map[string][]string{"organizationalUnit": {"orgunit1", "orgunit2"}},
// 							},
// 							Values: []string{"account1", "account2"},
// 						},
// 						"permissionSetArn": {
// 							Groupings: types.CreateAccessRuleTargetDetailArguments_Groupings{
// 								AdditionalProperties: map[string][]string{},
// 							},
// 							Values: []string{"abcdefg"},
// 						},
// 					},
// 				},
// 			},
// 			want: rule.Target{
// 				ProviderID:               "abcd",
// 				BuiltInProviderType:      "awssso",
// 				With:                     map[string]string{"permissionSetArn": "abcdefg"},
// 				WithSelectable:           map[string][]string{"accountId": {"account1", "account2"}},
// 				WithArgumentGroupOptions: map[string]map[string][]string{"accountId": {"organizationalUnit": {"orgunit1", "orgunit2"}}},
// 			},
// 			withProviderResponse: types.GetProviderResponse{
// 				HTTPResponse: &http.Response{
// 					StatusCode: http.StatusOK,
// 				},
// 				JSON200: &types.Provider{
// 					Id:   "abcd",
// 					Type: "awssso",
// 				},
// 			},
// 			withProviderArgsResponse: types.GetProviderArgsResponse{
// 				HTTPResponse: &http.Response{
// 					StatusCode: http.StatusOK,
// 				},
// 				JSON200: &ssov2Schema,
// 			},
// 		},
// 	}

// 	for _, tc := range testcases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			dbc := ddbmock.Client{
// 				PutErr: tc.wantErr,
// 			}

// 			clk := clock.NewMock()
// 			ctrl := gomock.NewController(t)

// 			defer ctrl.Finish()

// 			m := ahmocks.NewMockClientWithResponsesInterface(ctrl)

// 			m.EXPECT().GetProviderWithResponse(gomock.Any(), gomock.Eq(tc.give.ProviderId)).Return(&tc.withProviderResponse, nil)
// 			m.EXPECT().GetProviderArgsWithResponse(gomock.Any(), gomock.Eq(tc.give.ProviderId)).Return(&tc.withProviderArgsResponse, nil)

// 			cm := mocks.NewMockCacheService(ctrl)
// 			if !tc.dontExpectCacheCall {
// 				cm.EXPECT().RefreshCachedProviderArgOptions(gomock.Any(), gomock.Eq(tc.give.ProviderId), gomock.Any()).AnyTimes().Return(false, cacheArgOptionsResponse, cacheArgGroupOptionsResponse, nil)
// 			}
// 			s := Service{
// 				Clock:    clk,
// 				DB:       &dbc,
// 				AHClient: m,
// 				Cache:    cm,
// 			}
// 			got, err := s.ProcessTarget(context.Background(), tc.give, false)
// 			if tc.wantErr == nil {
// 				assert.NoError(t, err)
// 			} else {
// 				assert.EqualError(t, err, tc.wantErr.Error())
// 			}
// 			assert.Equal(t, tc.want, got)
// 		})
// 	}
// }
