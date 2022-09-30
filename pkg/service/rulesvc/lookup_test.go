package rulesvc

import (
	"context"
	"testing"

	"github.com/common-fate/ddb/ddbmock"
	"github.com/common-fate/granted-approvals/pkg/cache"
	"github.com/common-fate/granted-approvals/pkg/rule"
	"github.com/common-fate/granted-approvals/pkg/storage"
	"github.com/common-fate/granted-approvals/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestService_LookupRule(t *testing.T) {

	// we alias the rule to a variable to make test cases more concise to write.
	rule1AccountID := "123456789012"
	rule1PermissionSetARN := "arn:aws:sso:::permissionSet/ssoins-1234/ps-12341"

	rule1 := rule.AccessRule{
		ID: "test",
		Target: rule.Target{
			ProviderID:   "test-provider",
			ProviderType: "aws-sso",
			With: map[string]string{
				"accountId":        rule1AccountID,
				"permissionSetArn": rule1PermissionSetARN,
			},
		},
	}

	tests := []struct {
		name            string
		rules           []rule.AccessRule
		providerOptions []cache.ProviderOption
		args            LookupRuleOpts
		want            []LookedUpRule
		wantErr         bool
	}{
		{
			name:  "no rules returns nil not an error",
			rules: nil,
			args: LookupRuleOpts{
				ProviderType: "commonfate/aws-sso",
				Fields: LookupFields{
					AccountID: rule1AccountID,
					RoleName:  "GrantedAdministratorAccess",
				},
			},
			want: nil,
		},
		{
			name:  "single match",
			rules: []rule.AccessRule{rule1},
			providerOptions: []cache.ProviderOption{
				{
					Provider: "test-provider",
					Arg:      "permissionSetArn",
					Label:    "GrantedAdministratorAccess",
					Value:    rule1PermissionSetARN,
				},
			},
			args: LookupRuleOpts{
				ProviderType: "commonfate/aws-sso",
				Fields: LookupFields{
					AccountID: rule1AccountID,
					RoleName:  "GrantedAdministratorAccess",
				},
			},
			want: []LookedUpRule{
				{
					Rule: rule1,
				},
			},
		},
		{
			name: "multiple matches",
			rules: []rule.AccessRule{
				rule1,
				{
					ID: "second",
					Target: rule.Target{
						ProviderID:   "test-provider",
						ProviderType: "aws-sso",
						With: map[string]string{
							"accountId":        rule1AccountID,
							"permissionSetArn": rule1PermissionSetARN,
						},
					},
				},
			},
			providerOptions: []cache.ProviderOption{
				{
					Provider: "test-provider",
					Arg:      "permissionSetArn",
					Label:    "GrantedAdministratorAccess",
					Value:    rule1PermissionSetARN,
				},
			},
			args: LookupRuleOpts{
				ProviderType: "commonfate/aws-sso",
				Fields: LookupFields{
					AccountID: rule1AccountID,
					RoleName:  "GrantedAdministratorAccess",
				},
			},
			want: []LookedUpRule{
				{
					Rule: rule1,
				},
				{
					Rule: rule.AccessRule{
						ID: "second",
						Target: rule.Target{
							ProviderID:   "test-provider",
							ProviderType: "aws-sso",
							With: map[string]string{
								"accountId":        rule1AccountID,
								"permissionSetArn": rule1PermissionSetARN,
							},
						},
					},
				},
			},
		},
		{
			name:  "no match",
			rules: []rule.AccessRule{rule1},
			providerOptions: []cache.ProviderOption{
				{
					Provider: "test-provider",
					Arg:      "permissionSetArn",
					Label:    "GrantedAdministratorAccess",
					Value:    rule1PermissionSetARN,
				},
			},
			args: LookupRuleOpts{
				ProviderType: "commonfate/aws-sso",
				Fields: LookupFields{
					AccountID: rule1AccountID,
					// doesn't match the existing rules we have
					RoleName: "NoMatch",
				},
			},
			want: nil,
		},
		{
			name:  "match where permission set label has a description",
			rules: []rule.AccessRule{rule1},
			providerOptions: []cache.ProviderOption{
				{
					Provider: "test-provider",
					Arg:      "permissionSetArn",
					Label:    "GrantedAdministratorAccess: this is a test permission set",
					Value:    rule1PermissionSetARN,
				},
			},
			args: LookupRuleOpts{
				ProviderType: "commonfate/aws-sso",
				Fields: LookupFields{
					AccountID: rule1AccountID,
					RoleName:  "GrantedAdministratorAccess",
				},
			},
			want: []LookedUpRule{{Rule: rule1}},
		},
		{
			name:  "provider matches the permissions set arn by label but the rule does not contain the option should return no results",
			rules: []rule.AccessRule{rule1},
			providerOptions: []cache.ProviderOption{
				{
					Provider: "test-provider",
					Arg:      "permissionSetArn",
					Label:    "GrantedAdministratorAccess",
					Value:    "different to what the rule has",
				},
			},
			args: LookupRuleOpts{
				ProviderType: "commonfate/aws-sso",
				Fields: LookupFields{
					AccountID: rule1AccountID,
					RoleName:  "GrantedAdministratorAccess",
				},
			},
			want: nil,
		},
		{
			name: "selectable options are set",
			rules: []rule.AccessRule{
				{
					ID: "test",
					Target: rule.Target{
						ProviderID:   "test-provider",
						ProviderType: "aws-sso",
						WithSelectable: map[string][]string{
							"accountId":        {rule1AccountID, "second option"},
							"permissionSetArn": {rule1PermissionSetARN, "something else"},
						},
					},
				},
			},
			providerOptions: []cache.ProviderOption{
				{
					Provider: "test-provider",
					Arg:      "permissionSetArn",
					Label:    "GrantedAdministratorAccess",
					Value:    rule1PermissionSetARN,
				},
			},
			args: LookupRuleOpts{
				ProviderType: "commonfate/aws-sso",
				Fields: LookupFields{
					AccountID: rule1AccountID,
					RoleName:  "GrantedAdministratorAccess",
				},
			},
			want: []LookedUpRule{
				{
					Rule: rule.AccessRule{
						ID: "test",
						Target: rule.Target{
							ProviderID:   "test-provider",
							ProviderType: "aws-sso",
							WithSelectable: map[string][]string{
								"accountId":        {rule1AccountID, "second option"},
								"permissionSetArn": {rule1PermissionSetARN, "something else"},
							},
						},
					},
					SelectableWithOptionValues: []types.KeyValue{
						{
							Key:   "accountId",
							Value: rule1AccountID,
						},
						{
							Key:   "permissionSetArn",
							Value: rule1PermissionSetARN,
						},
					},
				},
			},
		},
		{
			name: "mix of matched rules with selectable and non-selectable options",
			rules: []rule.AccessRule{
				rule1,
				{
					ID: "test",
					Target: rule.Target{
						ProviderID:   "test-provider",
						ProviderType: "aws-sso",
						WithSelectable: map[string][]string{
							"accountId":        {rule1AccountID, "second option"},
							"permissionSetArn": {rule1PermissionSetARN, "something else"},
						},
					},
				},
			},
			providerOptions: []cache.ProviderOption{
				{
					Provider: "test-provider",
					Arg:      "permissionSetArn",
					Label:    "GrantedAdministratorAccess",
					Value:    rule1PermissionSetARN,
				},
			},
			args: LookupRuleOpts{
				ProviderType: "commonfate/aws-sso",
				Fields: LookupFields{
					AccountID: rule1AccountID,
					RoleName:  "GrantedAdministratorAccess",
				},
			},
			want: []LookedUpRule{
				{
					Rule: rule1,
				},
				{
					Rule: rule.AccessRule{
						ID: "test",
						Target: rule.Target{
							ProviderID:   "test-provider",
							ProviderType: "aws-sso",
							WithSelectable: map[string][]string{
								"accountId":        {rule1AccountID, "second option"},
								"permissionSetArn": {rule1PermissionSetARN, "something else"},
							},
						},
					},
					SelectableWithOptionValues: []types.KeyValue{
						{
							Key:   "accountId",
							Value: rule1AccountID,
						},
						{
							Key:   "permissionSetArn",
							Value: rule1PermissionSetARN,
						},
					},
				},
			},
		},
		{
			name: "no match if provider type is different",
			rules: []rule.AccessRule{
				{
					ID: "test",
					Target: rule.Target{
						ProviderID:   "test-provider",
						ProviderType: "something different",
						With: map[string]string{
							"accountId":        rule1AccountID,
							"permissionSetArn": rule1PermissionSetARN,
						},
					},
				},
			},
			providerOptions: []cache.ProviderOption{
				{
					Provider: "test-provider",
					Arg:      "permissionSetArn",
					Label:    "GrantedAdministratorAccess",
					Value:    "different to what the rule has",
				},
			},
			args: LookupRuleOpts{
				ProviderType: "commonfate/aws-sso",
				Fields: LookupFields{
					AccountID: rule1AccountID,
					RoleName:  "GrantedAdministratorAccess",
				},
			},
			want: nil,
		},
		{
			name:  "description label containing extra ':' characters",
			rules: []rule.AccessRule{rule1},
			providerOptions: []cache.ProviderOption{
				{
					Provider: "test-provider",
					Arg:      "permissionSetArn",
					// this should be an invalid description anyway, but test it just in case.
					Label: "GrantedAdministratorAccess: test::: test: : test",
					Value: rule1PermissionSetARN,
				},
			},
			args: LookupRuleOpts{
				ProviderType: "commonfate/aws-sso",
				Fields: LookupFields{
					AccountID: rule1AccountID,
					RoleName:  "GrantedAdministratorAccess",
				},
			},
			want: []LookedUpRule{{Rule: rule1}},
		},
		{
			name:  "no options for provider",
			rules: []rule.AccessRule{rule1},
			args: LookupRuleOpts{
				ProviderType: "commonfate/aws-sso",
				Fields: LookupFields{
					AccountID: rule1AccountID,
					RoleName:  "GrantedAdministratorAccess",
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := ddbmock.New(t)
			db.MockQuery(&storage.ListAccessRulesForGroupsAndStatus{Result: tt.rules})
			db.MockQuery(&storage.GetProviderOptions{Result: tt.providerOptions})

			s := &Service{
				DB: db,
			}
			got, err := s.LookupRule(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.LookupRule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
