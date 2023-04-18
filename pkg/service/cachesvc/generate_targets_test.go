package cachesvc

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/common-fate/common-fate/pkg/cache"
	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
)

func TestSync(t *testing.T) {
	type args struct {
		resources   []cache.TargetGroupResource
		accessRules []rule.AccessRule
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]map[string]Targets
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
					{ID: "accessRule_1", Targets: map[string]rule.Target{
						"targetgroup_1": {
							TargetGroupID: "targetgroup_1",
							Schema: providerregistrysdk.Target{
								Properties: map[string]providerregistrysdk.TargetField{
									"accountId":        {Resource: aws.String("Account")},
									"permissionSetArn": {Resource: aws.String("PermissionSet")},
								},
							},
						},
					}},
				},
			},
			want: map[string]map[string]Targets{
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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Sync() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOut(t *testing.T) {
	type args struct {
		in map[string]map[string]Targets
	}
	tests := []struct {
		name string
		args args
		want map[string]Target
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
			},
			want: map[string]Target{
				"targetgroup_1#accountId#account_1#permissionSetArn#permissionSet_1": {
					fields: map[string]string{
						"accountId":        "account_1",
						"permissionSetArn": "permissionSet_1",
					},
					rules: []string{"accessRule_1", "accessRule_2"},
				},
				"targetgroup_1#accountId#account_1#permissionSetArn#permissionSet_2": {
					fields: map[string]string{
						"accountId":        "account_1",
						"permissionSetArn": "permissionSet_2",
					},
					rules: []string{"accessRule_1"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateDistinctTargets(tt.args.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Out() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterForRules(t *testing.T) {
	type args struct {
		targets []Target
		rules   []string
	}

	t1 := Target{
		fields: map[string]string{"hello": "world"},
		rules:  []string{"accessRule_1"},
	}
	t2 := Target{
		fields: map[string]string{"hello": "world"},
		rules:  []string{"accessRule_2"},
	}
	tests := []struct {
		name string
		args args
		want []Target
	}{
		{
			name: "ok",
			args: args{
				targets: []Target{t1, t2},
				rules:   []string{"accessRule_1"},
			},
			want: []Target{t1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tf := NewTargetFilter(tt.args.rules)
			tf.Filter(tt.args.targets)
			if got := tf.Dump(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilterForRules() = %v, want %v", got, tt.want)
			}
		})
	}
}
