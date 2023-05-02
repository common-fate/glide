package cache

import (
	"reflect"
	"testing"
)

func TestFilterTargetsByAccessRules(t *testing.T) {
	type args struct {
		targets []Target
		rules   []string
	}

	t1 := Target{
		Fields: []Field{{
			ID:    "hello",
			Value: "world",
		}},
		AccessRules: map[string]AccessRule{"accessRule_1": {}},
	}
	t2 := Target{
		Fields: []Field{{
			ID:    "hello",
			Value: "world",
		}},
		AccessRules: map[string]AccessRule{"accessRule_2": {}},
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
			tf := NewFilterTargetsByAccessRule(tt.args.rules)
			tf.Filter(tt.args.targets)
			if got := tf.Dump(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilterForRules() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestFilterTargetsByGroups(t *testing.T) {
	type args struct {
		targets []Target
		groups  []string
	}

	t1 := Target{
		Fields: []Field{{
			ID:    "hello",
			Value: "world",
		}},
		IDPGroupsWithAccess: map[string]struct{}{"group_1": {}},
	}
	t2 := Target{
		Fields: []Field{{
			ID:    "hello",
			Value: "world",
		}},
		IDPGroupsWithAccess: map[string]struct{}{"group_2": {}},
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
				groups:  []string{"group_1"},
			},
			want: []Target{t1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tf := NewFilterTargetsByGroups(tt.args.groups)
			tf.Filter(tt.args.targets)
			if got := tf.Dump(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilterForRules() = %v, want %v", got, tt.want)
			}
		})
	}
}
