package cachesvc

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/common-fate/common-fate/pkg/cache"
	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/common-fate/common-fate/pkg/target"
	"github.com/stretchr/testify/assert"
)

func TestSync(t *testing.T) {
	type args struct {
		resources   []cache.TargetGroupResource
		accessRules []rule.AccessRule
	}
	tests := []struct {
		name    string
		args    args
		want    resourceAccessRuleMapping
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				resources: []cache.TargetGroupResource{
					{TargetGroupID: "targetgroup_1", ResourceType: "Account", Resource: cache.Resource{ID: "account_1", Name: "Account_1"}},
					{TargetGroupID: "targetgroup_1", ResourceType: "PermissionSet", Resource: cache.Resource{ID: "permissionSet_1", Name: "PermissionSet_1"}},
					{TargetGroupID: "targetgroup_1", ResourceType: "PermissionSet", Resource: cache.Resource{ID: "permissionSet_2", Name: "PermissionSet_2"}},
				},
				accessRules: []rule.AccessRule{
					{ID: "accessRule_1", Targets: []rule.Target{
						{
							TargetGroup: target.Group{ID: "targetgroup_1", Schema: target.GroupSchema{
								Target: target.TargetSchema{
									Properties: map[string]target.TargetField{
										"accountId":        {Resource: aws.String("Account")},
										"permissionSetArn": {Resource: aws.String("PermissionSet")},
									},
								},
							}},
						},
					}},
				},
			},
			want: resourceAccessRuleMapping{
				"accessRule_1": {
					"targetgroup_1": Targets{
						map[string]cache.Resource{
							"accountId":        {ID: "account_1", Name: "Account_1"},
							"permissionSetArn": {ID: "permissionSet_1", Name: "PermissionSet_1"},
						},
						map[string]cache.Resource{
							"accountId":        {ID: "account_1", Name: "Account_1"},
							"permissionSetArn": {ID: "permissionSet_2", Name: "PermissionSet_2"},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createResourceAccessRuleMapping(tt.args.resources, tt.args.accessRules)

			if (err != nil) != tt.wantErr {
				t.Errorf("Sync() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGenerateDistinctTargets(t *testing.T) {
	type args struct {
		in         map[string]map[string]Targets
		acessRules []rule.AccessRule
	}
	tests := []struct {
		name string
		args args
		want []cache.Target
	}{
		{
			name: "ok",
			args: args{
				in: map[string]map[string]Targets{
					"accessRule_1": {
						"targetgroup_1": Targets{
							map[string]cache.Resource{
								"accountId":        {ID: "account_1", Name: "Account_1"},
								"permissionSetArn": {ID: "permissionSet_1", Name: "PermissionSet_1"},
							},
							map[string]cache.Resource{
								"accountId":        {ID: "account_1", Name: "Account_1"},
								"permissionSetArn": {ID: "permissionSet_2", Name: "PermissionSet_2"},
							},
						},
					},
					"accessRule_2": {
						"targetgroup_1": Targets{
							map[string]cache.Resource{
								"accountId":        {ID: "account_1", Name: "Account_1"},
								"permissionSetArn": {ID: "permissionSet_1", Name: "PermissionSet_1"},
							},
						},
					},
				},
				acessRules: []rule.AccessRule{
					{ID: "accessRule_1", Groups: []string{"group_1", "group_2"},
						Targets: []rule.Target{{TargetGroup: target.Group{ID: "targetgroup_1", Schema: target.GroupSchema{
							Target: target.TargetSchema{
								Properties: map[string]target.TargetField{
									"accountId":        {Title: aws.String("AWS Account")},
									"permissionSetArn": {Title: aws.String("AWS Permission Set"), Description: aws.String("a permission set field description")},
								},
							},
						},
						}}}},
					{ID: "accessRule_2", Groups: []string{"group_3", "group_4"},
						Targets: []rule.Target{{TargetGroup: target.Group{ID: "targetgroup_1", Schema: target.GroupSchema{
							Target: target.TargetSchema{
								Properties: map[string]target.TargetField{
									"accountId":        {Title: aws.String("AWS Account")},
									"permissionSetArn": {Title: aws.String("AWS Permission Set"), Description: aws.String("a permission set field description")},
								},
							},
						}}}}},
				},
			},
			want: []cache.Target{
				{

					Fields: []cache.Field{
						{
							ID:         "accountId",
							Value:      "account_1",
							ValueLabel: "Account_1",
							FieldTitle: "AWS Account",
						},
						{
							ID:               "permissionSetArn",
							Value:            "permissionSet_1",
							FieldTitle:       "AWS Permission Set",
							FieldDescription: aws.String("a permission set field description"),
							ValueLabel:       "PermissionSet_1",
						},
					},

					IDPGroupsWithAccess: map[string]struct{}{"group_1": {}, "group_2": {}, "group_3": {}, "group_4": {}},
					AccessRules: map[string]cache.AccessRule{
						"accessRule_2": {MatchedTargetGroups: []string{"targetgroup_1"}},
						"accessRule_1": {MatchedTargetGroups: []string{"targetgroup_1"}}},
				},
				{

					Fields: []cache.Field{
						{
							ID:         "accountId",
							Value:      "account_1",
							ValueLabel: "Account_1",
							FieldTitle: "AWS Account",
						},
						{
							ID:               "permissionSetArn",
							Value:            "permissionSet_2",
							FieldTitle:       "AWS Permission Set",
							FieldDescription: aws.String("a permission set field description"),
							ValueLabel:       "PermissionSet_2",
						},
					},
					IDPGroupsWithAccess: map[string]struct{}{"group_1": {}, "group_2": {}},
					AccessRules:         map[string]cache.AccessRule{"accessRule_1": {MatchedTargetGroups: []string{"targetgroup_1"}}},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generateDistinctTargets(tt.args.in, tt.args.acessRules)
			assert.Equal(t, tt.want, got)
		})
	}
}
