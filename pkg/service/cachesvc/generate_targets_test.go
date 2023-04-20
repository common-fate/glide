package cachesvc

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/common-fate/common-fate/pkg/cache"
	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/common-fate/common-fate/pkg/target"
	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
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
							TargetGroup: target.Group{ID: "targetgroup_1", Schema: providerregistrysdk.Target{
								Properties: map[string]providerregistrysdk.TargetField{
									"accountId":        {Resource: aws.String("Account")},
									"permissionSetArn": {Resource: aws.String("PermissionSet")},
								},
							}},
						},
					}},
				},
			},
			want: resourceAccessRuleMapping{
				"accessRule_1": {
					"targetgroup_1": Targets{
						map[string]string{
							"accountId":        "account_1",
							"permissionSetArn": "permissionSet_1",
						},
						map[string]string{
							"accountId":        "account_1",
							"permissionSetArn": "permissionSet_2",
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
							map[string]string{
								"accountId":        "account_1",
								"permissionSetArn": "permissionSet_1",
							},
							map[string]string{
								"accountId":        "account_1",
								"permissionSetArn": "permissionSet_2",
							},
						},
					},
					"accessRule_2": {
						"targetgroup_1": Targets{
							map[string]string{
								"accountId":        "account_1",
								"permissionSetArn": "permissionSet_1",
							},
						},
					},
				},
				acessRules: []rule.AccessRule{
					{ID: "accessRule_1", Groups: []string{"group_1", "group_2"}},
					{ID: "accessRule_2", Groups: []string{"group_3", "group_4"}},
				},
			},
			want: []cache.Target{
				{
					TargetGroupID: "targetgroup_1",
					Fields: []cache.Field{
						{
							ID:    "accountId",
							Value: "account_1",
						},
						{
							ID:    "permissionSetArn",
							Value: "permissionSet_1",
						},
					},

					Groups:      map[string]struct{}{"group_1": {}, "group_2": {}, "group_3": {}, "group_4": {}},
					AccessRules: map[string]struct{}{"accessRule_2": {}, "accessRule_1": {}},
				},
				{
					TargetGroupID: "targetgroup_1",
					Fields: []cache.Field{
						{
							ID:    "accountId",
							Value: "account_1",
						},
						{
							ID:    "permissionSetArn",
							Value: "permissionSet_2",
						},
					},
					Groups:      map[string]struct{}{"group_1": {}, "group_2": {}},
					AccessRules: map[string]struct{}{"accessRule_1": {}},
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
