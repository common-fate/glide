package preflightsvc

import (
	"context"
	"reflect"
	"testing"

	"github.com/benbjohnson/clock"
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/cache"
	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb/ddbmock"
	"github.com/stretchr/testify/assert"
)

func TestGroupTargets(t *testing.T) {
	clk := clock.NewMock()
	target1 := cache.Target{

		Kind: cache.Kind{
			Publisher: "publisher1",
			Name:      "target1",
			Kind:      "kind1",
		},
		AccessRules: map[string]cache.AccessRule{
			"rule1": {},
		},
		Groups: map[string]struct{}{
			"group1": {},
			"group2": {},
		},
		Fields: []cache.Field{
			{
				ID:         "field1",
				FieldTitle: "Field 1",
				ValueLabel: "Value 1",
				Value:      "value1",
			},
			{
				ID:         "field2",
				FieldTitle: "Field 2",
				ValueLabel: "Value 2",
				Value:      "value2",
			},
		},
	}

	target2 := cache.Target{
		Kind: cache.Kind{
			Publisher: "publisher2",
			Name:      "target1",
			Kind:      "kind1",
		},
		AccessRules: map[string]cache.AccessRule{
			"rule2": {},
		},
		Groups: map[string]struct{}{
			"group1": {},
			"group2": {},
		},
		Fields: []cache.Field{
			{
				ID:         "field1",
				FieldTitle: "Field 1",
				ValueLabel: "Value 1",
				Value:      "value1",
			},
			{
				ID:         "field2",
				FieldTitle: "Field 2",
				ValueLabel: "Value 2",
				Value:      "value2",
			},
		},
	}

	tests := []struct {
		name               string
		targets            []cache.Target
		AccessGroups       []access.PreflightAccessGroup
		wantErr            bool
		mockGetAccessRule1 rule.AccessRule
		mockGetAccessRule2 rule.AccessRule
	}{

		{
			name:    "multiple targets with diff access rules creates multiple groups",
			targets: []cache.Target{target1, target2},
			AccessGroups: []access.PreflightAccessGroup{
				{
					Targets: []cache.Target{
						target1,
					},
					TimeConstraints: types.AccessRuleTimeConstraints{
						MaxDurationSeconds: 3600,
					},
				},
				{
					Targets: []cache.Target{
						target2,
					},
					TimeConstraints: types.AccessRuleTimeConstraints{
						MaxDurationSeconds: 3600,
					},
				},
			},
			mockGetAccessRule1: rule.AccessRule{
				ID:          "rule1",
				Description: "string",
				Name:        "string",
				Groups:      []string{"string"},

				Approval: rule.Approval{
					Groups: []string{"a"},
					Users:  []string{"b"},
				},
				TimeConstraints: types.AccessRuleTimeConstraints{MaxDurationSeconds: 3600},
			},
			mockGetAccessRule2: rule.AccessRule{
				ID:          "rule2",
				Description: "string",
				Name:        "string",
				Groups:      []string{"string"},

				Approval: rule.Approval{
					Groups: []string{"a"},
					Users:  []string{"b"},
				},
				TimeConstraints: types.AccessRuleTimeConstraints{MaxDurationSeconds: 3600},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := ddbmock.New(t)

			db.MockQueries(
				&storage.GetAccessRule{Result: &tt.mockGetAccessRule1},
				&storage.GetAccessRule{Result: &tt.mockGetAccessRule2},
			)

			s := &Service{
				DB:    db,
				Clock: clk,
			}

			got, _ := s.GroupTargets(context.Background(), tt.targets)

			//override ids
			for i, _ := range tt.AccessGroups {
				tt.AccessGroups[i].ID = got[i].ID
			}

			assert.Equal(t, tt.AccessGroups, got)
		})
	}
}

func TestGetAccessRuleWithLongerDuration(t *testing.T) {
	// Define two AccessRule objects with different configurations
	ar1 := rule.AccessRule{
		ID: "test",

		Approval:        rule.Approval{Groups: []string{"admin"}},
		TimeConstraints: types.AccessRuleTimeConstraints{MaxDurationSeconds: 60},
	}
	ar2 := rule.AccessRule{
		ID: "test",

		TimeConstraints: types.AccessRuleTimeConstraints{MaxDurationSeconds: 120},
	}

	// Call the function and verify the result matches the expected AccessRule
	expectedResult := ar2
	result := CompareAccessRules(ar1, ar2)
	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("GetAccessRuleWithLongerDuration() failed. Expected result: %v, but got: %v", expectedResult, result)
	}

	// Define two other AccessRule objects with different configurations
	ar3 := rule.AccessRule{
		ID: "test",

		Approval:        rule.Approval{Users: []string{"user1"}},
		TimeConstraints: types.AccessRuleTimeConstraints{MaxDurationSeconds: 180},
	}
	ar4 := rule.AccessRule{
		ID: "test",

		Approval:        rule.Approval{Users: []string{"user2"}},
		TimeConstraints: types.AccessRuleTimeConstraints{MaxDurationSeconds: 240},
	}

	// Call the function and verify the result matches the expected AccessRule
	expectedResult = ar4
	result = CompareAccessRules(ar3, ar4)
	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("GetAccessRuleWithLongerDuration() failed. Expected result: %v, but got: %v", expectedResult, result)
	}

	// Define two AccessRule objects with the same configuration
	ar5 := rule.AccessRule{
		ID: "test",

		Approval:        rule.Approval{Groups: []string{"admin"}},
		TimeConstraints: types.AccessRuleTimeConstraints{MaxDurationSeconds: 60},
	}
	ar6 := rule.AccessRule{
		ID: "test",

		Approval:        rule.Approval{Groups: []string{"admin"}},
		TimeConstraints: types.AccessRuleTimeConstraints{MaxDurationSeconds: 60},
	}

	// Call the function and verify the result matches one of the input AccessRule objects
	result = CompareAccessRules(ar5, ar6)
	if !reflect.DeepEqual(result, ar5) && !reflect.DeepEqual(result, ar6) {
		t.Errorf("GetAccessRuleWithLongerDuration() failed. Expected result to be one of the input AccessRules, but got: %v", result)
	}

	// Define two AccessRule objects with the same configuration but different MaxDurationSeconds values
	ar7 := rule.AccessRule{
		ID: "test",

		Approval:        rule.Approval{Groups: []string{"admin"}},
		TimeConstraints: types.AccessRuleTimeConstraints{MaxDurationSeconds: 120},
	}
	ar8 := rule.AccessRule{
		ID: "test",

		Approval:        rule.Approval{Groups: []string{"admin"}},
		TimeConstraints: types.AccessRuleTimeConstraints{MaxDurationSeconds: 60},
	}

	// Call the function and verify the result matches the AccessRule with the longer MaxDurationSeconds value
	expectedResult = ar7
	result = CompareAccessRules(ar7, ar8)
	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("GetAccessRuleWithLongerDuration() failed. Expected result: %v, but got: %v", expectedResult, result)
	}

	// Define two AccessRule objects with the same configuration but different MaxDurationSeconds values
	ar9 := rule.AccessRule{
		ID:              "test",
		Approval:        rule.Approval{},
		TimeConstraints: types.AccessRuleTimeConstraints{MaxDurationSeconds: 120},
	}
	ar10 := rule.AccessRule{
		ID: "test",

		Approval:        rule.Approval{},
		TimeConstraints: types.AccessRuleTimeConstraints{MaxDurationSeconds: 60},
	}

	// Call the function and verify the result matches the AccessRule with the longer MaxDurationSeconds value
	expectedResult = ar9
	result = CompareAccessRules(ar9, ar10)
	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("GetAccessRuleWithLongerDuration() failed. Expected result: %v, but got: %v", expectedResult, result)
	}

	ar11 := rule.AccessRule{
		Approval:        rule.Approval{},
		TimeConstraints: types.AccessRuleTimeConstraints{MaxDurationSeconds: 120},
	}
	ar12 := rule.AccessRule{
		ID: "test",

		Approval:        rule.Approval{},
		TimeConstraints: types.AccessRuleTimeConstraints{MaxDurationSeconds: 60},
	}

	// Call the function and verify the result matches the AccessRule with the longer MaxDurationSeconds value
	expectedResult = ar12
	result = CompareAccessRules(ar11, ar12)
	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("GetAccessRuleWithLongerDuration() failed. Expected result: %v, but got: %v", expectedResult, result)
	}

	ar13 := rule.AccessRule{
		ID: "test",

		Approval:        rule.Approval{},
		TimeConstraints: types.AccessRuleTimeConstraints{MaxDurationSeconds: 120},
	}
	ar14 := rule.AccessRule{

		Approval:        rule.Approval{},
		TimeConstraints: types.AccessRuleTimeConstraints{MaxDurationSeconds: 60},
	}

	// Call the function and verify the result matches the AccessRule with the longer MaxDurationSeconds value
	expectedResult = ar13
	result = CompareAccessRules(ar13, ar14)
	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("GetAccessRuleWithLongerDuration() failed. Expected result: %v, but got: %v", expectedResult, result)
	}

}

func TestValidateNoDuplicates(t *testing.T) {

	tests := []struct {
		name    string
		args    types.CreatePreflightRequest
		wantErr bool
	}{
		{
			name: "No duplicates",
			args: types.CreatePreflightRequest{
				Targets: []string{"target1", "target2", "target3"},
			},
			wantErr: false,
		},
		{
			name: "Duplicate targets",
			args: types.CreatePreflightRequest{
				Targets: []string{"target1", "target2", "target1"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateNoDuplicates(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("ValidateNoDuplicates() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_ValidateAccessToAllTargets(t *testing.T) {
	clk := clock.NewMock()
	type args struct {
		user             identity.User
		preflightRequest types.CreatePreflightRequest
	}
	tests := []struct {
		name          string
		args          args
		mockGetTarget cache.Target
		want          []cache.Target
		wantErr       bool
	}{
		{
			name: "ok",
			args: args{
				user: identity.User{Groups: []string{"group_1"}},
				preflightRequest: types.CreatePreflightRequest{
					Targets: []string{"tg_1"},
				},
			},
			mockGetTarget: cache.Target{
				Groups: map[string]struct{}{"group_1": {}},
			},
			want: []cache.Target{{
				Groups: map[string]struct{}{"group_1": {}},
			}},
		},
		{
			name: "fail",
			args: args{
				user: identity.User{Groups: []string{"group_1"}},
				preflightRequest: types.CreatePreflightRequest{
					Targets: []string{"tg_1"},
				},
			},
			mockGetTarget: cache.Target{
				Groups: map[string]struct{}{"group_2": {}},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := ddbmock.New(t)

			db.MockQuery(&storage.GetCachedTarget{
				Result: &tt.mockGetTarget,
			})
			s := &Service{
				DB:    db,
				Clock: clk,
			}
			got, err := s.ValidateAccessToAllTargets(context.Background(), tt.args.user, tt.args.preflightRequest)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.ValidateAccessToAllTargets() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.ValidateAccessToAllTargets() = %v, want %v", got, tt.want)
			}
		})
	}
}
