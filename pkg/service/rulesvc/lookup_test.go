package rulesvc

// import (
// 	"context"
// 	"reflect"
// 	"testing"

// 	"github.com/common-fate/common-fate/pkg/cache"
// 	"github.com/common-fate/common-fate/pkg/identity"
// 	"github.com/common-fate/common-fate/pkg/rule"
// 	"github.com/common-fate/common-fate/pkg/storage"
// 	"github.com/common-fate/common-fate/pkg/types"
// 	"github.com/common-fate/ddb/ddbmock"
// 	"github.com/stretchr/testify/assert"
// )

// func TestService_LookupRule(t *testing.T) {

// 	// we alias the rule to a variable to make test cases more concise to write.
// 	rule1AccountID := "123456789012"
// 	rule1PermissionSetARN := "arn:aws:sso:::permissionSet/ssoins-1234/ps-12341"

// 	// fill with "G1", "G2"
// 	groups := []string{"g1"}

// 	rule1 := rule.AccessRule{
// 		ID: "test",
// 		Target: rule.Target{
// 			TargetGroupID: "test-provider",

// 			With: map[string]string{
// 				"accountId":        rule1AccountID,
// 				"permissionSetArn": rule1PermissionSetARN,
// 			},
// 		},
// 		Groups: groups,
// 	}

// 	tests := []struct {
// 		name                   string
// 		rules                  []rule.AccessRule
// 		providerOptions        []cache.ProviderOption
// 		providerArgGroupOption *cache.ProviderArgGroupOption
// 		args                   LookupRuleOpts
// 		want                   []LookedUpRule
// 		wantErr                bool
// 	}{
// 		{
// 			name:  "no rules returns nil not an error",
// 			rules: nil,
// 			args: LookupRuleOpts{
// 				ProviderType: "commonfate/aws-sso",
// 				Fields: LookupFields{
// 					AccountID: rule1AccountID,
// 					RoleName:  "GrantedAdministratorAccess",
// 				},
// 				User: identity.User{
// 					Groups: rule1.Groups,
// 				},
// 			},
// 			want: nil,
// 		},
// 		{
// 			name:  "single match",
// 			rules: []rule.AccessRule{rule1},
// 			providerOptions: []cache.ProviderOption{
// 				{
// 					Provider: "test-provider",
// 					Arg:      "permissionSetArn",
// 					Label:    "GrantedAdministratorAccess",
// 					Value:    rule1PermissionSetARN,
// 				},
// 			},
// 			args: LookupRuleOpts{
// 				ProviderType: "commonfate/aws-sso",
// 				Fields: LookupFields{
// 					AccountID: rule1AccountID,
// 					RoleName:  "GrantedAdministratorAccess",
// 				},
// 				User: identity.User{
// 					Groups: rule1.Groups,
// 				},
// 			},
// 			want: []LookedUpRule{
// 				{
// 					Rule: rule1,
// 				},
// 			},
// 		},
// 		{
// 			name: "multiple matches",
// 			rules: []rule.AccessRule{
// 				rule1,
// 				{
// 					ID: "second",
// 					Target: rule.Target{
// 						TargetGroupID: "test-provider",

// 						With: map[string]string{
// 							"accountId":        rule1AccountID,
// 							"permissionSetArn": rule1PermissionSetARN,
// 						},
// 					},
// 					Groups: rule1.Groups,
// 				},
// 			},
// 			providerOptions: []cache.ProviderOption{
// 				{
// 					Provider: "test-provider",
// 					Arg:      "permissionSetArn",
// 					Label:    "GrantedAdministratorAccess",
// 					Value:    rule1PermissionSetARN,
// 				},
// 			},
// 			args: LookupRuleOpts{
// 				ProviderType: "commonfate/aws-sso",
// 				Fields: LookupFields{
// 					AccountID: rule1AccountID,
// 					RoleName:  "GrantedAdministratorAccess",
// 				},
// 				User: identity.User{
// 					Groups: rule1.Groups,
// 				},
// 			},
// 			want: []LookedUpRule{
// 				{
// 					Rule: rule1,
// 				},
// 				{
// 					Rule: rule.AccessRule{
// 						ID: "second",
// 						Target: rule.Target{
// 							TargetGroupID: "test-provider",
// 							With: map[string]string{
// 								"accountId":        rule1AccountID,
// 								"permissionSetArn": rule1PermissionSetARN,
// 							},
// 						},
// 						Groups: rule1.Groups,
// 					},
// 				},
// 			},
// 		},
// 		{
// 			name:  "no match",
// 			rules: []rule.AccessRule{rule1},
// 			providerOptions: []cache.ProviderOption{
// 				{
// 					Provider: "test-provider",
// 					Arg:      "permissionSetArn",
// 					Label:    "GrantedAdministratorAccess",
// 					Value:    rule1PermissionSetARN,
// 				},
// 			},
// 			args: LookupRuleOpts{
// 				ProviderType: "commonfate/aws-sso",
// 				Fields: LookupFields{
// 					AccountID: rule1AccountID,
// 					// doesn't match the existing rules we have
// 					RoleName: "NoMatch",
// 				},
// 				User: identity.User{
// 					Groups: rule1.Groups,
// 				},
// 			},
// 			want: nil,
// 		},

// 		{
// 			name:  "provider matches the permissions set arn by label but the rule does not contain the option should return no results",
// 			rules: []rule.AccessRule{rule1},
// 			providerOptions: []cache.ProviderOption{
// 				{
// 					Provider: "test-provider",
// 					Arg:      "permissionSetArn",
// 					Label:    "GrantedAdministratorAccess",
// 					Value:    "different to what the rule has",
// 				},
// 			},
// 			args: LookupRuleOpts{
// 				ProviderType: "commonfate/aws-sso",
// 				Fields: LookupFields{
// 					AccountID: rule1AccountID,
// 					RoleName:  "GrantedAdministratorAccess",
// 				},
// 				User: identity.User{
// 					Groups: rule1.Groups,
// 				},
// 			},
// 			want: nil,
// 		},
// 		{
// 			name: "selectable options are set",
// 			rules: []rule.AccessRule{
// 				{
// 					ID: "test",
// 					Target: rule.Target{
// 						TargetGroupID: "test-provider",
// 						// WithSelectable: map[string][]string{
// 						// 	"accountId":        {rule1AccountID, "second option"},
// 						// 	"permissionSetArn": {rule1PermissionSetARN, "something else"},
// 						// },
// 					},
// 					Groups: rule1.Groups,
// 				},
// 			},
// 			providerOptions: []cache.ProviderOption{
// 				{
// 					Provider: "test-provider",
// 					Arg:      "permissionSetArn",
// 					Label:    "GrantedAdministratorAccess",
// 					Value:    rule1PermissionSetARN,
// 				},
// 			},
// 			args: LookupRuleOpts{
// 				ProviderType: "commonfate/aws-sso",
// 				Fields: LookupFields{
// 					AccountID: rule1AccountID,
// 					RoleName:  "GrantedAdministratorAccess",
// 				},
// 				User: identity.User{
// 					Groups: rule1.Groups,
// 				},
// 			},
// 			want: []LookedUpRule{
// 				{
// 					Rule: rule.AccessRule{
// 						ID: "test",
// 						Target: rule.Target{
// 							TargetGroupID: "test-provider",
// 							// WithSelectable: map[string][]string{
// 							// 	"accountId":        {rule1AccountID, "second option"},
// 							// 	"permissionSetArn": {rule1PermissionSetARN, "something else"},
// 							// },
// 						},
// 						Groups: rule1.Groups,
// 					},
// 					SelectableWithOptionValues: []types.KeyValue{
// 						{
// 							Key:   "accountId",
// 							Value: rule1AccountID,
// 						},
// 						{
// 							Key:   "permissionSetArn",
// 							Value: rule1PermissionSetARN,
// 						},
// 					},
// 				},
// 			},
// 		},
// 		{
// 			name: "mix of matched rules with selectable and non-selectable options",
// 			rules: []rule.AccessRule{
// 				rule1,
// 				{
// 					ID: "test",
// 					Target: rule.Target{
// 						TargetGroupID: "test-provider",
// 						// WithSelectable: map[string][]string{
// 						// 	"accountId":        {rule1AccountID, "second option"},
// 						// 	"permissionSetArn": {rule1PermissionSetARN, "something else"},
// 						// },
// 					},
// 					Groups: rule1.Groups,
// 				},
// 			},
// 			providerOptions: []cache.ProviderOption{
// 				{
// 					Provider: "test-provider",
// 					Arg:      "permissionSetArn",
// 					Label:    "GrantedAdministratorAccess",
// 					Value:    rule1PermissionSetARN,
// 				},
// 			},
// 			args: LookupRuleOpts{
// 				ProviderType: "commonfate/aws-sso",
// 				Fields: LookupFields{
// 					AccountID: rule1AccountID,
// 					RoleName:  "GrantedAdministratorAccess",
// 				},
// 				User: identity.User{
// 					Groups: rule1.Groups,
// 				},
// 			},
// 			want: []LookedUpRule{
// 				{
// 					Rule: rule1,
// 				},
// 				{
// 					Rule: rule.AccessRule{
// 						ID: "test",
// 						Target: rule.Target{
// 							TargetGroupID: "test-provider",
// 							// WithSelectable: map[string][]string{
// 							// 	"accountId":        {rule1AccountID, "second option"},
// 							// 	"permissionSetArn": {rule1PermissionSetARN, "something else"},
// 							// },
// 						},
// 						Groups: rule1.Groups,
// 					},
// 					SelectableWithOptionValues: []types.KeyValue{
// 						{
// 							Key:   "accountId",
// 							Value: rule1AccountID,
// 						},
// 						{
// 							Key:   "permissionSetArn",
// 							Value: rule1PermissionSetARN,
// 						},
// 					},
// 				},
// 			},
// 		},
// 		{
// 			name: "no match if provider type is different",
// 			rules: []rule.AccessRule{
// 				{
// 					ID: "test",
// 					Target: rule.Target{
// 						TargetGroupID: "test-provider",
// 						With: map[string]string{
// 							"accountId":        rule1AccountID,
// 							"permissionSetArn": rule1PermissionSetARN,
// 						},
// 					},
// 					Groups: rule1.Groups,
// 				},
// 			},
// 			providerOptions: []cache.ProviderOption{
// 				{
// 					Provider: "test-provider",
// 					Arg:      "permissionSetArn",
// 					Label:    "GrantedAdministratorAccess",
// 					Value:    "different to what the rule has",
// 				},
// 			},
// 			args: LookupRuleOpts{
// 				ProviderType: "commonfate/aws-sso",
// 				Fields: LookupFields{
// 					AccountID: rule1AccountID,
// 					RoleName:  "GrantedAdministratorAccess",
// 				},
// 				User: identity.User{
// 					Groups: rule1.Groups,
// 				},
// 			},
// 			want: nil,
// 		},
// 		{
// 			name:  "no options for provider",
// 			rules: []rule.AccessRule{rule1},
// 			args: LookupRuleOpts{
// 				ProviderType: "commonfate/aws-sso",
// 				Fields: LookupFields{
// 					AccountID: rule1AccountID,
// 					RoleName:  "GrantedAdministratorAccess",
// 				},
// 				User: identity.User{
// 					Groups: rule1.Groups,
// 				},
// 			},
// 			want: nil,
// 		},
// 		{
// 			name: "match via arg group",
// 			rules: []rule.AccessRule{
// 				{
// 					ID: "test",
// 					Target: rule.Target{
// 						TargetGroupID: "test-provider",
// 						With: map[string]string{
// 							"permissionSetArn": rule1PermissionSetARN,
// 						},
// 						// WithArgumentGroupOptions: map[string]map[string][]string{
// 						// 	"accountId": {
// 						// 		"organizationalUnit": {"orgUnit1"},
// 						// 	},
// 						// },
// 					},
// 					Groups: rule1.Groups,
// 				},
// 			},
// 			providerOptions: []cache.ProviderOption{
// 				{
// 					Provider: "test-provider",
// 					Arg:      "permissionSetArn",
// 					Label:    "GrantedAdministratorAccess",
// 					Value:    rule1PermissionSetARN,
// 				},
// 			},
// 			args: LookupRuleOpts{
// 				ProviderType: "commonfate/aws-sso",
// 				Fields: LookupFields{
// 					AccountID: rule1AccountID,
// 					RoleName:  "GrantedAdministratorAccess",
// 				},
// 				User: identity.User{
// 					Groups: rule1.Groups,
// 				},
// 			},
// 			providerArgGroupOption: &cache.ProviderArgGroupOption{
// 				Provider: "test-provider",
// 				Arg:      "accountId",
// 				Group:    "organizationalUnit",
// 				Value:    "orgUnit1",
// 				Children: []string{rule1AccountID},
// 			},
// 			want: []LookedUpRule{
// 				{
// 					Rule: rule.AccessRule{
// 						ID: "test",
// 						Target: rule.Target{
// 							TargetGroupID: "test-provider",
// 							With: map[string]string{
// 								"permissionSetArn": rule1PermissionSetARN,
// 							},
// 							// WithArgumentGroupOptions: map[string]map[string][]string{
// 							// 	"accountId": {
// 							// 		"organizationalUnit": {"orgUnit1"},
// 							// 	},
// 							// },
// 						},
// 						Groups: rule1.Groups,
// 					},
// 				},
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			db := ddbmock.New(t)
// 			db.MockQuery(&storage.ListAccessRulesForStatus{Result: tt.rules})
// 			db.MockQuery(&storage.ListCachedProviderOptionsForArg{Result: tt.providerOptions})
// 			db.MockQuery(&storage.GetCachedProviderArgGroupOptionValueForArg{Result: tt.providerArgGroupOption})

// 			s := &Service{
// 				DB: db,
// 			}
// 			got, err := s.LookupRule(context.Background(), tt.args)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("Service.LookupRule() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			assert.Equal(t, tt.want, got)
// 		})
// 	}
// }

// func Test_FilterRulesByGroupMap(t *testing.T) {
// 	type args struct {
// 		groups []string
// 		rules  []rule.AccessRule
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want []rule.AccessRule
// 	}{
// 		{
// 			name: "no groups",
// 			args: args{groups: []string{}},
// 			want: []rule.AccessRule{},
// 		},
// 		{
// 			name: "no rules",
// 			args: args{groups: []string{"group1"}},
// 			want: []rule.AccessRule{},
// 		},
// 		{
// 			name: "no matches",
// 			args: args{
// 				groups: []string{"group1"},
// 				rules: []rule.AccessRule{
// 					{
// 						ID:     "rule1",
// 						Target: rule.Target{},
// 						Groups: []string{"group2"},
// 					},
// 				},
// 			},
// 			want: []rule.AccessRule{},
// 		},
// 		{
// 			name: "match",
// 			args: args{
// 				groups: []string{"group2"},
// 				rules: []rule.AccessRule{
// 					{
// 						ID:     "rule1",
// 						Target: rule.Target{},
// 						Groups: []string{"group2"},
// 					},
// 				},
// 			},
// 			want: []rule.AccessRule{
// 				{
// 					ID:     "rule1",
// 					Target: rule.Target{},
// 					Groups: []string{"group2"},
// 				},
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := FilterRulesByGroupMap(tt.args.groups, tt.args.rules); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("filterRulesByGroupMap() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
