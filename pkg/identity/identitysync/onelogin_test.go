package identitysync

import (
	"reflect"
	"testing"

	"github.com/common-fate/common-fate/pkg/identity"
)

func Test_filterGroups(t *testing.T) {
	type args struct {
		groups       []identity.IDPGroup
		filterString string
	}
	tests := []struct {
		name    string
		args    args
		want    []identity.IDPGroup
		wantErr bool
	}{
		{
			name: "regex starts with: single match",
			args: args{
				groups: []identity.IDPGroup{
					{Name: "foo"},
					{Name: "bar"},
				},
				filterString: "^foo",
			},
			want: []identity.IDPGroup{
				{Name: "foo"},
			},
			wantErr: false,
		},
		{
			name: "regex ends with: double match",
			args: args{
				groups: []identity.IDPGroup{
					{Name: "foo"},
					{Name: "oooooo"},
					{Name: "bar"},
				},
				filterString: "oo$",
			},
			want: []identity.IDPGroup{
				{Name: "foo"},
				{Name: "oooooo"},
			},
			wantErr: false,
		},
		{
			name: "regex starts with: empty reponse",
			args: args{
				groups: []identity.IDPGroup{
					{Name: "bar"},
				},
				filterString: "^foo",
			},
			want:    []identity.IDPGroup{},
			wantErr: false,
		},
		{
			name: "invalid regex throws err",
			args: args{
				groups: []identity.IDPGroup{
					{Name: "bar"},
				},
				filterString: "^(foo",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FilterGroups(tt.args.groups, tt.args.filterString)
			if (err != nil) != tt.wantErr {
				t.Errorf("filterGroups() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterGroups() = %v, want %v", got, tt.want)
			}
		})
	}
}
