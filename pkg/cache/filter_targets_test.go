package cache

import (
	"reflect"
	"testing"
)

func TestFilterForRules(t *testing.T) {
	type args struct {
		targets []Target
		rules   []string
	}

	t1 := Target{
		Fields:      map[string]string{"hello": "world"},
		AccessRules: []string{"accessRule_1"},
	}
	t2 := Target{
		Fields:      map[string]string{"hello": "world"},
		AccessRules: []string{"accessRule_2"},
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
